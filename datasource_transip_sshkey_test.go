package main

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"time"

	"testing"
)

func TestAccTransipDataSourceSSHKey(t *testing.T) {
	timestamp := time.Now().Unix()

	testFixture := fmt.Sprintf(`
	resource "transip_sshkey" "test" {
		description = "test-%d"
    key         = "%s"
	}
	`, timestamp, testSSHKey)

	testConfig := fmt.Sprintf(`
	data "transip_sshkey" "test" {
		description = "test-%d"
	}
	`, timestamp)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testFixture,
			},
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.transip_sshkey.test", "key", testSSHKey),
				),
			},
		},
	})
}
