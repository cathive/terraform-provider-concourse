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

	flyrc_path, err := flyRcLocation()
	if err != nil {
		return err
	}

	flyrc_contents, err := flyReadConfig(flyrc_path)
	if err != nil {
		return err
	}

	yaml.Unmarshal(*flyrc_contents, cfg)

	return nil
}

func flyRcLocation() (*string, error) {
	// Check if an ENV var has been set with a path
	// Todo: Find out if this is the correct ENV var, or if it fly even has one.
	if flyrc, ok := os.LookupEnv("FLYRC"); ok {
		return &flyrc, nil
	}

	// Otherwise, return the default flyrc location
	cu, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("unable to determine current user for reading the .flyrc file: %v", err)
	}
	flyrc := fmt.Sprintf("%s/.flyrc", cu.HomeDir)
	return &flyrc, nil
}

// Get the bytes of the flyrc config based on the filepath given
func flyReadConfig(flyrc *string) (*[]byte, error) {
	if _, err := os.Stat(*flyrc); err != nil {
		return nil, fmt.Errorf("unable to stat the flyrc file (%s): %v", *flyrc, err)
	}

	config_bytes, err := ioutil.ReadFile(*flyrc)

	return &config_bytes, err
}
