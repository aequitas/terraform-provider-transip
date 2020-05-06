package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTransipDataSourcePrivateNetwork(t *testing.T) {
	var testConfig = `data "transip_private_network" "test" {name = "%s"}`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testConfig, os.Getenv("TF_VAR_private_network_name")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.transip_private_network.test", "name", os.Getenv("TF_VAR_private_network_name")),
				),
			},
		},
	})
}
