package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataTeam() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTeamRead,
		Exists: resourceTeamExists,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Team name",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}
