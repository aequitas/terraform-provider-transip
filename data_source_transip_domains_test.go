package main

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTransipDataSourceDomains(t *testing.T) {
	if v := os.Getenv("TF_VAR_domain"); v == "" {
		t.Skip("TF_VAR_domain must be set for acceptance tests")
	}

	var testConfig = `data "transip_domains" "all_domains" { }`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.transip_domains", "all_domains"),
				),
			},
		},
	})
}
