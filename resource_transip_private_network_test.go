package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"

	"os"
	"testing"

	"github.com/hashicorp/terraform/terraform"
)

func TestAccTransipResourcePrivateNetwork(t *testing.T) {
	testConfig := fmt.Sprintf(`
	resource "transip_private_network" "test" {
    	description = "%s"
	}
	`, os.Getenv("TF_VAR_private_network_name"))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:        testConfig,
				ResourceName:  "transip_private_network.test",
				ImportState:   true,
				ImportStateId: os.Getenv("TF_VAR_private_network_name"),
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if s[0].ID != os.Getenv("TF_VAR_private_network_name") {
						return fmt.Errorf("import failed")
					}
					return nil
				},
			},
		},
	})
}
