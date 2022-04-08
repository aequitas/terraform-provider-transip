package main

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/sethvargo/go-password/password"

	"os"
	"testing"
)

func TestAccTransipResourcOpenstackUser(t *testing.T) {
	if os.Getenv("THIS_IS_GOING_TO_COST_ME_MONEY") == "" {
		t.Skip("THIS_IS_GOING_TO_COST_ME_MONEY not set, skipping")
	}
	openstackProjectId := os.Getenv("TF_VAR_openstack_project_id")
	if openstackProjectId == "" {
		t.Skip("TF_VAR_openstack_project_id not provided, skipping")
	}

	pwd, err := password.Generate(8, 1, 1, false, false)
	if err != nil {
		t.Skip("failed to generate password", err)
	}

	testConfig := fmt.Sprintf(`
	resource "transip_openstack_user" "test" {
		projectId = "%s"
		username = "tf-test-user"
		email = "tf-test-user@transip.nl"
		password = "%s"
		description = "terraform test user"
	}
	`, openstackProjectId, pwd)

	testConfigUpdate := fmt.Sprintf(`
	resource "transip_openstack_user" "test" {
		projectId = "%s"
		username = "tf-test-user"
		email = "tf-test-user-updated-email@transip.nl"
		password = "%s"
		description = "terraform test user"
	}
	`, openstackProjectId, pwd)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("transip_openstack_user", "username", "tf-test-user"),
				),
			},
			{
				Config: testConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("transip_openstack_user", "email", "tf-test-user-updated-email@transip.nl"),
				),
			},
		},
	})
}
