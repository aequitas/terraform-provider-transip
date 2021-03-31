package main

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTransipResourceVpsImport(t *testing.T) {
	vpsName := os.Getenv("TF_VAR_vps_name")
	if vpsName == "" {
		t.Skip("TF_VAR_vps_name not provided, skipping")
	}

	testConfig := fmt.Sprintf(`
	resource "transip_vps" "test" {
    product_name     = "vps-bladevps-x1"
    operating_system = "Debian 6"
	}
	`)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:        testConfig,
				ResourceName:  "transip_vps.test",
				ImportState:   true,
				ImportStateId: vpsName,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if s[0].ID != vpsName {
						return fmt.Errorf("import failed")
					}
					return nil
				},
			},
		},
	})
}

func TestAccTransipResourceVps(t *testing.T) {
	if os.Getenv("THIS_IS_GOING_TO_COST_ME_MONEY") == "" {
		t.Skip("THIS_IS_GOING_TO_COST_ME_MONEY not set, skipping")
	}
	timeStamp := time.Now().Unix()
	testConfig := fmt.Sprintf(`
	resource "transip_vps" "test" {
		description             = "test-%d"
    product_name     = "vps-bladevps-x1"
    operating_system = "ubuntu-20.04"
	}
	`, timeStamp)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("transip_vps.test", "status", "running"),
				),
			},
		},
	})
}

func TestAccTransipResourceVpsUpdate(t *testing.T) {
	if os.Getenv("THIS_IS_GOING_TO_COST_ME_MONEY") == "" {
		t.Skip("THIS_IS_GOING_TO_COST_ME_MONEY not set, skipping")
	}
	timestamp := time.Now().Unix()
	// TODO: Add tags when supported.
	testConfig := fmt.Sprintf(`
	resource "transip_vps" "test" {
		description             = "test-%d"
	product_name     = "vps-bladevps-x1"
	operating_system = "ubuntu-20.04"
	is_customer_locked = true
	}
	`, timestamp)

	testConfig2 := fmt.Sprintf(`
	resource "transip_vps" "test" {
		description             = "test-%d-updated"
	product_name     = "vps-bladevps-x1"
	operating_system = "ubuntu-20.04"
	is_customer_locked = false
	}
	`, timestamp)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					// Test if lock is set and if the description has changed.
					resource.TestCheckResourceAttr("transip_vps.test", "is_customer_locked", "true"),
					resource.TestCheckResourceAttr("transip_vps.test", "description", fmt.Sprintf("test-%d", timestamp)),
				),
			},
			{
				Config: testConfig2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("transip_vps.test", "is_customer_locked", "false"),
					resource.TestCheckResourceAttr("transip_vps.test", "description", fmt.Sprintf("test-%d-updated", timestamp)),
				),
			},
		},
	})
}
