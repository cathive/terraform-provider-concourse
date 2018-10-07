package concourse

import (
	"testing"
	"github.com/hashicorp/terraform/terraform"
	"github.com/hashicorp/terraform/helper/schema"
)

func testAccPreCheck(t *testing.T) {}
var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init()  {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"concourse": testAccProvider,
	}
}