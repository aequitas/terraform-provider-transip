package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTransipDataSourcePrivateNetwork(t *testing.T) {
	privateNetworkName := os.Getenv("TF_VAR_private_network_name")
	if privateNetworkName == "" {
		t.Skip("TF_VAR_private_network_name not provided, skipping")
	}
	var testConfig = `data "transip_private_network" "test" {name = "%s"}`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testConfig, privateNetworkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.transip_private_network.test", "name", privateNetworkName),
				),
			},
		},
	})
}
