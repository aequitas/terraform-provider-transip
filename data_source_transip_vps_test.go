package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTransipDataSourceVps(t *testing.T) {
	var testConfig = `data "transip_vps" "test" {name = "%s"}`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testConfig, os.Getenv("TRANSIP_TEST_VPS")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.transip_vps.test", "status", "running"),
				),
			},
		},
	})
}
