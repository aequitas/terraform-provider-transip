package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTransipDataSourceDomain(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					// TODO: does TestCheckResourceAttrSet not work on lists?
					// resource.TestCheckResourceAttrSet("data.transip_domain.test", "nameservers"),
					resource.TestCheckResourceAttrSet("data.transip_domain.test", "is_locked"),
				),
			},
		},
	})
}

var testConfig = `
data "transip_domain" "test" {
	name = "locohost.nl"
}
`
