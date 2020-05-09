package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"

	"os"
	"testing"

	"github.com/hashicorp/terraform/terraform"
)

func TestAccTransipResourcePrivateNetwork(t *testing.T) {
	privateNetworkName := os.Getenv("TF_VAR_private_network_name")
	if privateNetworkName == "" {
		t.Skip("TF_VAR_private_network_name not provided, skipping")
	}
	testConfig := fmt.Sprintf(`
	resource "transip_private_network" "test" {
    	description = "%s"
	}
	`, privateNetworkName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:        testConfig,
				ResourceName:  "transip_private_network.test",
				ImportState:   true,
				ImportStateId: privateNetworkName,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if s[0].ID != privateNetworkName {
						return fmt.Errorf("import failed")
					}
					return nil
				},
			},
		},
	})
}
