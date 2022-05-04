package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTransipDataSourceVps(t *testing.T) {
	vpsName := os.Getenv("TF_VAR_vps_name")
	if vpsName == "" {
		t.Skip("TF_VAR_vps_name not provided, skipping")
	}
	var testConfig = `data "transip_vps" "test" {name = "%s"}`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testConfig, vpsName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.transip_vps.test", "status", "running"),
				),
			},
		},
	})
}
