package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTransipDataSourceDomain(t *testing.T) {
	if v := os.Getenv("TF_VAR_domain"); v == "" {
		t.Fatal("TF_VAR_domain must be set for acceptance tests")
	}

	var testConfig = `data "transip_domain" "test" {name = "%s"}`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testConfig, os.Getenv("TF_VAR_domain")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.transip_domain.test", "is_transfer_locked", "false"),
				),
			},
		},
	})
}
