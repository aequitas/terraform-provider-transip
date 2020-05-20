package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTransipResourceVpsFirewall(t *testing.T) {
	vpsName := os.Getenv("TF_VAR_vps_name")
	testConfig := fmt.Sprintf(`

	resource "transip_vps_firewall" "test" {
	  name = "%s"

	  inbound_rule {
	    description = "HTTP"
	    port        = "80"
	    protocol    = "tcp"
	  }

	  inbound_rule {
	    description = "HTTPS"
	    port        = 443
	    protocol    = "udp"
	  }

	  inbound_rule {
	    description = "SSH"
	    port        = 22
	    protocol    = "tcp_udp"
	  }

	  inbound_rule {
	    description = "RANGE"
	    port        = "90-100"
	    protocol    = "tcp"
	  }
	}
	`, vpsName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:        testConfig,
				ResourceName:  "transip_vps_firewall.test",
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
