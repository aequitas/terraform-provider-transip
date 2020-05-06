package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"os"
	"testing"
)

func TestAccTransipResourcePrivateNetworkAttachment(t *testing.T) {
	testConfig := fmt.Sprintf(`
	resource "transip_private_network_attachment" "test" {
		private_network_id = "%s"
		vps_id = "%s"
	}
	`, os.Getenv("TF_VAR_private_network_name"), os.Getenv("TF_VAR_vps_name"))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:        testConfig,
				ResourceName:  "transip_private_network_attachment.test",
				ImportState:   false,
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
