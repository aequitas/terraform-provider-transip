package main

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"os"
	"testing"
)

func TestAccTransipResourcePrivateNetworkAttachment(t *testing.T) {
	privateNetworkName := os.Getenv("TF_VAR_private_network_name")
	if privateNetworkName == "" {
		t.Skip("TF_VAR_private_network_name not provided, skipping")
	}
	vpsName := os.Getenv("TF_VAR_vps_name")
	if vpsName == "" {
		t.Skip("TF_VAR_vps_name not provided, skipping")
	}

	testConfig := fmt.Sprintf(`
	resource "transip_private_network_attachment" "test" {
		private_network_id = "%s"
		vps_id = "%s"
	}
	`, privateNetworkName, os.Getenv("TF_VAR_vps_name"))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:        testConfig,
				ResourceName:  "transip_private_network_attachment.test",
				ImportState:   false,
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
