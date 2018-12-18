package main

import (
	"github.com/cathive/terraform-provider-concourse/concourse"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return concourse.Provider()
		},
	})
}
