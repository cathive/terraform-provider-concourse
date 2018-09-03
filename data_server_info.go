package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataServerInfoRead(d *schema.ResourceData, m interface{}) error {
	cfg := m.(Config)

	d.SetId(cfg.Concourse().URL())

	if err := d.Set("version", cfg.Version()); err != nil {
		return fmt.Errorf("unable to set version field: %v", err)
	}

	if err := d.Set("worker_version", cfg.WorkerVersion()); err != nil {
		return fmt.Errorf("unable to set worker_version field: %v", err)
	}

	return nil

}

func dataServerInfoExists(d *schema.ResourceData, m interface{}) (bool, error) {
	// Well,... since we fetch the server info during initialization of the provider
	// itself (if initialization fails, we throw an error there!), we won't have
	// to do anything at this point...
	return true, nil
}

func dataServerInfo() *schema.Resource {
	return &schema.Resource{
		Read:   dataServerInfoRead,
		Exists: dataServerInfoExists,
		Schema: map[string]*schema.Schema{
			"version": {
				Description: "Concourse ATC/web server version",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"worker_version": {
				Description: "Worker version",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}
