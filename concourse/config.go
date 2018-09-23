package concourse

import (
	"encoding/json"
	"fmt"
	"github.com/concourse/go-concourse/concourse"
	"net/http"
		"net/url"
)

// Config provides access to all the stuff that is necessary to properly operate
// upon the Concourse ATC.
type Config interface {
	Concourse() concourse.Client
	Version() string
	WorkerVersion() string
	UserInfo() *SkyUserInfo
}

type config struct {
	url           string
	insecure      bool
	team          string
	client        concourse.Client
	userInfo      *SkyUserInfo
	version       string
	workerVersion string
}

func (c *config) Concourse() concourse.Client {
	return c.client
}

func (c *config) Version() string {
	return c.version
}

func (c *config) WorkerVersion() string {
	return c.workerVersion
}

func (c *config) UserInfo() *SkyUserInfo {
	return c.userInfo
}

// NewConfig creates a new configuration structure to be used provider-internally
// During initialization of the config structure, we will send 2 requests to the
// Concourse CI instance. First of all, we figure out the version of the ATC system
// and secondly, we fetch the user information from the sky marshal user-info endpoint.
func NewConfig(url *url.URL, httpClient *http.Client, insecure bool, team string) (Config, error) {

	client := concourse.NewClient(url.String(), httpClient, false)

	info, err := client.GetInfo()
	if err != nil {
		return nil, fmt.Errorf("unable to contact Concourse CI: %v", err)
	}

	userInfoURL := fmt.Sprintf("%s/%s", client.URL(), "sky/userinfo")
	resp, err := client.HTTPClient().Get(userInfoURL)
	if err != nil {
		return nil, fmt.Errorf("unable to communicate with the Concourse CI API server: %v", err)
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("user is not authorized to communuicate with the Concourse CI API server (%s returned status code %d)", userInfoURL, http.StatusUnauthorized)
	}

	userInfo := &SkyUserInfo{}
	if err := json.NewDecoder(resp.Body).Decode(userInfo); err != nil {
		return nil, fmt.Errorf("unable to gather user information: %v", err)
	}

	return &config{
		url:           url.String(),
		insecure:      insecure,
		team:          team,
		client:        client,
		userInfo:      userInfo,
		version:       info.Version,
		workerVersion: info.WorkerVersion,
	}, nil

}
