package concourse

import (
	"fmt"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"os"
	"os/user"
)

// FlyRC is a representation of the configuration file structure that is stored by the
// "fly" command line interface. (Usually to be found in ~/.flyrc)
type FlyRC struct {
	Filename string
	Targets map[string]Target `yaml:"targets"`
}

type Target struct {
	API      string `yaml:"api"`
	Team     string `yaml:"team"`
	Insecure bool   `yaml:"insecure,omitempty"`
	Token    struct {
		Type  string `yaml:"type"`
		Value string `yaml:"value"`
	} `yaml:"token"`
}

// Reads in a `flyrc` file and returns a FlyRC struct
func (rc *FlyRC) ImportConfig() error {
	cfg := FlyRC{}

	rc.setFlyRcLocation()

	flyrc_contents, err := rc.readFlyConfig()
	if err != nil {
		return err
	}

	yaml.Unmarshal(*flyrc_contents, cfg)

	return nil
}

func (rc *FlyRC) setFlyRcLocation() {
	fallback := ".flyrc" // If all else fails, we'll just return .flyrc in the current directory

	// Check if an ENV var has been set with a path
	// Todo: Find out if this is the correct ENV var, or if it fly even has one.
	if flyrc, ok := os.LookupEnv("FLYRC"); ok {
		rc.Filename = flyrc
	}

	// Otherwise, return the default flyrc location
	cu, err := user.Current()
	if err != nil {
		rc.Filename = fallback
	}
	flyrc := fmt.Sprintf("%s/.flyrc", cu.HomeDir)
	rc.Filename = flyrc
}

// Get the bytes of the flyrc config based on the filepath given
func (rc *FlyRC) readFlyConfig() (*[]byte, error) {
	if _, err := os.Stat(rc.Filename); err != nil {
		return nil, fmt.Errorf("unable to stat the flyrc file (%s): %v", rc.Filename, err)
	}

	config_bytes, err := ioutil.ReadFile(rc.Filename)

	return &config_bytes, err
}
