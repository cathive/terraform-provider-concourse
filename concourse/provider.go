package concourse

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"gopkg.in/yaml.v2"
	"net/http"
		"net/url"
	"os"
	"os/user"
	"strings"
	"golang.org/x/oauth2"
)

// FlyRC is a representation of the configuration file structure that is stored by the
// "fly" command line interface. (Usually to be found in ~/.flyrc)
type FlyRC struct {
	Targets map[string]struct {
		API      string `yaml:"api"`
		Team     string `yaml:"team"`
		Insecure bool   `yaml:"insecure,omitempty"`
		Token    struct {
			Type  string `yaml:"type"`
			Value string `yaml:"value"`
		} `yaml:"token"`
	} `yaml:"targets"`
}

// SkyUserInfo encapsulates all the information that is being reported by the Sky marshal
// "sky/userinfo" REST endpoint
type SkyUserInfo struct {
	CSRF     string   `json:"csrf"`
	Email    string   `json:"email"`
	Exp      int      `json:"exp"`
	IsAdmin  bool     `json:"is_admin"`
	Name     string   `json:"name"`
	Sub      string   `json:"sub"`
	Teams    []string `json:"teams"`
	UserID   string   `json:"user_id"`
	UserName string   `json:"user_name"`
}

// Provider creates a new Concourse terraform provider instance.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"concourse_url": {
				Description:   "Concourse URL to authenticate with",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"target"},
			},
			"insecure": {
				Description:   "Skip verification of the endpoint's SSL certificate",
				Type:          schema.TypeBool,
				Default:       false,
				Optional:      true,
				ConflictsWith: []string{"target"},
			},
			"auth_token_type": {
				Description:   "Authentication token type",
				Type:          schema.TypeString,
				Optional:      true,
				Default: "Bearer",
				ConflictsWith: []string{"target"},
			},
			"auth_token_value": {
				Description: "Authentication token value",
				Type: schema.TypeString,
				Optional: true,
				ConflictsWith: []string{"target"},
			},
			"target": {
				Description:   "ID of the concourse target if NOT using any of the other parameters",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"concourse_url", "insecure", "auth_token_type", "auth_token_value"},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"concourse_team": resourceTeam(),
			"concourse_pipeline": resourcePipeline(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"concourse_caller_identity": dataCallerIdentity(),
			"concourse_server_info":     dataServerInfo(),
			"concourse_team":            dataTeam(),
		},
		ConfigureFunc: configure,
	}
}

func configure(d *schema.ResourceData) (interface{}, error) {
	concourseURL := d.Get("concourse_url").(string)
	insecure := d.Get("insecure").(bool)
	authTokenType := d.Get("auth_token_type").(string)
	authTokenValue := d.Get("auth_token_value").(string)
	targetName := d.Get("target").(string)

	var u *url.URL

	// Let's try to read the fly CLI configuration file if the user did not specify
	// any connection parameters in the provider configuration.
	if targetName != "" {
		cu, err := user.Current()
		if err != nil {
			return nil, fmt.Errorf("unable to determine current user: %v", err)
		}
		cfgFilePath := fmt.Sprintf("%s/.flyrc", cu.HomeDir)
		if _, err := os.Stat(cfgFilePath); err != nil {
			return nil, fmt.Errorf("unable to find Fly configuration file (%s): %v", cfgFilePath, err)
		}
		cfgFile, err := os.Open(cfgFilePath)
		if err != nil {
			return nil, fmt.Errorf("unable to open Fly configuration file (%s): %v", cfgFile.Name(), err)
		}
		cfg := FlyRC{}
		if yaml.NewDecoder(cfgFile).Decode(&cfg); err != nil {
			return nil, fmt.Errorf("unable to parse Fly configuration file (%s): %v", cfgFile.Name(), err)
		}
		if len(cfg.Targets) <= 0 {
			return nil, fmt.Errorf("no targets found in Fly configuration file (%s)", cfgFile.Name())
		}
		if targetName == "" {
			return nil, fmt.Errorf("provider argument \"target\" must be specified")
		}
		target, exists := cfg.Targets[targetName]
		if !exists {
			return nil, fmt.Errorf("unable to find targetName with ID \"%s\" in Fly configuration file %s", targetName, cfgFile.Name())
		}
		concourseURL = target.API
		if u, err = url.Parse(concourseURL); err != nil {
			return nil, fmt.Errorf("unable to parse URL (%s): %v", concourseURL, err)
		}
		insecure = target.Insecure
		authTokenType = target.Token.Type
		authTokenValue = target.Token.Value
	} else {
		cfgMissing := make([]string, 0)
		if concourseURL == "" {
			cfgMissing = append(cfgMissing, "\"concourse_url\"")
		}
		if authTokenType == "" {
			cfgMissing = append(cfgMissing, "\"auth_token_type\"")
		}
		if authTokenValue == "" {
			cfgMissing = append(cfgMissing, "\"auth_token_value\"")
		}
		if len(cfgMissing) > 0 {
			return nil, fmt.Errorf("required configuration parameter(s) missing: %s", strings.Join(cfgMissing, ", "))
		}
		u, err := url.Parse(concourseURL)
		if err != nil {
			return nil, fmt.Errorf("unable to parse URL (%s): %v", concourseURL, err)
		}

		u = &url.URL{
			Scheme: u.Scheme,
			Host:   u.Host,
			Path:   u.Path,
		}

	}
	oAuthToken := &oauth2.Token{
		TokenType: authTokenType,
		AccessToken: authTokenValue,
	}
	transport := &oauth2.Transport{
		Source: oauth2.StaticTokenSource(oAuthToken),
	}
	httpClient := &http.Client{
		Transport: transport,
	}

	return NewConfig(u, httpClient, insecure, targetName)
}
