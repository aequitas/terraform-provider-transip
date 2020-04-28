package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"

	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

func TestAccTransipResourceVps(t *testing.T) {
	vpsName := os.Getenv("TF_VAR_vps_name")
	testConfig := fmt.Sprintf(`
	resource "transip_vps" "test" {
		name             = "%s"
    product_name     = "vps-bladevps-x1"
    operating_system = "Debian 6"
	}
	`, vpsName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:        testConfig,
				ResourceName:  "transip_vps.test",
				ImportState:   true,
				ImportStateId: vpsName,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if s[0].ID != vpsName {
						return fmt.Errorf("import failed")
					}
					return nil
				},
			},
		},
	})
}
