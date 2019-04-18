package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTransipDataSourceDomain(t *testing.T) {
	var testConfig = `data "transip_domain" "test" {name = "%s"}`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testConfig, os.Getenv("TRANSIP_TEST_DOMAIN")),
				Check: resource.ComposeTestCheckFunc(
					// TODO: does TestCheckResourceAttrSet not work on lists?
					// resource.TestCheckResourceAttrSet("data.transip_domain.test", "nameservers"),
					resource.TestCheckResourceAttrSet("data.transip_domain.test", "is_locked"),
				),
			},
		},
	})
}
