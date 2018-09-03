package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataCallerIdentityRead(d *schema.ResourceData, m interface{}) error {
	userInfo := m.(Config).UserInfo()

	d.SetId(userInfo.UserID)

	if err := d.Set("user_name", userInfo.UserName); err != nil {
		return fmt.Errorf("unable to set user_name field: %v", err)
	}

	if err := d.Set("is_admin", userInfo.IsAdmin); err != nil {
		return fmt.Errorf("unable to set is_admin field: %v", err)
	}

	return nil

}

func dataCallerIdentityExists(d *schema.ResourceData, m interface{}) (bool, error) {
	// Well,... since we already use the caller identity during initialization of the provider
	// itself (if initialization fails, we throw an error there!), we won't have
	// to do anything at this point...
	return true, nil
}

func dataCallerIdentity() *schema.Resource {
	return &schema.Resource{
		Read:   dataCallerIdentityRead,
		Exists: dataCallerIdentityExists,
		Schema: map[string]*schema.Schema{
			"user_name": {
				Description: "User name of the current Concourse ATC API user",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"is_admin": {
				Description: "Flag that indicates if the current Concourse ATC API user has admin permissions",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}
