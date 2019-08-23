package concourse

import (
	"fmt"
	"github.com/concourse/concourse/fly/rc"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"strings"
)

// Provider creates a new Concourse terraform provider instance.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"concourse_url": {
				Description:   "Concourse URL to authenticate with",
				Type:          schema.TypeString,
				Optional:      true,
			},
			"insecure": {
				Description:   "Skip verification of the endpoint's SSL certificate",
				Type:          schema.TypeBool,
				Default:       false,
				Optional:      true,
			},
			"auth_token_type": {
				Description:   "Authentication token type",
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "Bearer",
			},
			"auth_token_value": {
				Description:   "Authentication token value",
				Type:          schema.TypeString,
				Optional:      true,
			},
			"target": {
				Description:   "ID of the concourse target if NOT using any of the other parameters",
				Type:          schema.TypeString,
				Optional:      true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"concourse_team":     resourceTeam(),
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
		target, err := rc.LoadTarget(rc.TargetName(targetName), false)
		if err != nil {
			return nil, fmt.Errorf("unable to load target from flyrc: %v", err)
		}

		concourseURL = target.URL()
		if u, err = url.Parse(concourseURL); err != nil {
			return nil, fmt.Errorf("unable to parse URL (%s): %v", concourseURL, err)
		}
		insecure = target.TLSConfig().InsecureSkipVerify
		authTokenType = target.Token().Type
		authTokenValue = target.Token().Value
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
		TokenType:   authTokenType,
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
