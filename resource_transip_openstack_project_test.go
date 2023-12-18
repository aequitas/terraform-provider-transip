package main

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"os"
	"testing"
)

func TestAccTransipResourcOpenstackProjectImport(t *testing.T) {
	openstackProjectName := os.Getenv("TF_VAR_openstack_project_name")
	if openstackProjectName == "" {
		t.Skip("TF_VAR_openstack_project_name not provided, skipping")
	}

	testConfig := fmt.Sprintf(`
	resource "transip_openstack_project" "test" {
		name = "%s"
		description = "terraform test project"
	}
	`, openstackProjectName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:       testConfig,
				ResourceName: "transip_openstack_project.test",
				ImportState:  false,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if s[0].ID != openstackProjectName {
						return fmt.Errorf("import failed")
					}
					return nil
				},
			},
		},
	})
}

func TestAccTransipResourcOpenstackProject(t *testing.T) {
	if os.Getenv("THIS_IS_GOING_TO_COST_ME_MONEY") == "" {
		t.Skip("THIS_IS_GOING_TO_COST_ME_MONEY not set, skipping")
	}

	timestamp := time.Now().Unix()
	testConfig := `
	resource "transip_openstack_project" "test" {
		name = "aequitasterraformtest-tf-test"
		description = "terraform test project"
	}
	`

	testConfigUpdate := fmt.Sprintf(`
	resource "transip_openstack_project" "test" {
		name = "tf-test-%d"
		description = "terraform test project %d"
	}
	`, timestamp, timestamp)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("transip_openstack_project", "name", "tf-test"),
				),
			},
			{
				Config: testConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("transip_openstack_project", "name", fmt.Sprintf("tf-test-%d", timestamp)),
				),
			},
		},
	})
}
