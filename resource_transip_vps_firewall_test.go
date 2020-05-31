package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/transip/gotransip/v6/repository"
	"github.com/transip/gotransip/v6/vps"
)

func TestAccTransipResourceVpsFirewall(t *testing.T) {
	vpsName := os.Getenv("TF_VAR_vps_name")
	if vpsName == "" {
		t.Skip("TF_VAR_vps_name not provided, skipping")
	}

	var firewall vps.Firewall

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTransipResourceVpsFirewall(vpsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransipResourceVpsFirewallExists("transip_vps_firewall.test", &firewall),
					testAccCheckVpsFirewallAttributes(&firewall),
					resource.TestCheckResourceAttr("transip_vps_firewall.test", "is_enabled", "true"),
					resource.TestCheckResourceAttr("transip_vps_firewall.test", "name", vpsName),
					resource.TestCheckResourceAttr("transip_vps_firewall.test", "inbound_rule.#", "6"),
				),
			},
		},
	})
}

// Check if all remote attributes match up
func testAccCheckVpsFirewallAttributes(firewall *vps.Firewall) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(firewall.RuleSet) != 6 {
			return fmt.Errorf("firewall does not contain 6 inbound rules")
		}

		return nil
	}
}

// Query the API to verify remote resource exist
func testAccCheckTransipResourceVpsFirewallExists(n string, firewall *vps.Firewall) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := testAccProvider.Meta().(repository.Client)
		repository := vps.FirewallRepository{Client: client}
		v, err := repository.GetFirewall(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to obtain vps firewall %q: %s", rs.Primary.ID, err)
		}

		// If no error, assign the response
		*firewall = v
		return nil
	}
}

// returns an configuration for a firewall given a specific vps (name)
func testAccTransipResourceVpsFirewall(name string) string {
	return fmt.Sprintf(`

	resource "transip_vps_firewall" "test" {
	  name = "%s"

	  inbound_rule {
	    description = "HTTP"
	    port        = "80"
	    protocol    = "tcp"
	    whitelist   = [
	        "192.168.0.2/32",
	        "::/128",
	        "127.0.0.2/32",
	        "127.0.0.1/32",
	        "192.168.0.1/32",
	    ]
	  }

	  inbound_rule {
	    description = "HTTPS"
	    port        = 443
	    protocol    = "tcp"

	  }

	  inbound_rule {
	    description = "SSH"
	    port        = 22
	    protocol    = "tcp"
	    whitelist   = []
	  }

	  inbound_rule {
	    description = "RANGE"
	    port        = "90-100"
	    protocol    = "tcp"
	  }

	  inbound_rule {
	    description = "UDP"
	    port        = "200"
	    protocol    = "udp"
	  }

	  inbound_rule {
	    description = "Both"
	    port        = "300"
	    protocol    = "tcp_udp"
	  }
	}
	`, name)
}
