package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTransipDataSourceOpenstackUser(t *testing.T) {
	openstackUserName := os.Getenv("TF_VAR_openstack_username")
	if openstackUserName == "" {
		t.Skip("TF_VAR_openstack_username not provided, skipping")
	}
	var testConfig = `data "transip_openstack_user" "test" {username = "%s"}`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testConfig, openstackUserName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.transip_openstack_username.test", "username", openstackUserName),
				),
			},
		},
	})
}
