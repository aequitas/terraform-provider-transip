package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTransipDataSourceOpenstackProject(t *testing.T) {
	openstackProjectName := os.Getenv("TF_VAR_openstack_project_name")
	if openstackProjectName == "" {
		t.Skip("TF_VAR_openstack_project_name not provided, skipping")
	}
	var testConfig = `data "transip_openstack_project" "test" {name = "%s"}`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testConfig, openstackProjectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.transip_openstack_project.test", "name", openstackProjectName),
				),
			},
		},
	})
}
