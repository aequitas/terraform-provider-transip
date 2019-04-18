package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTransipResourceDomain(t *testing.T) {
	timestamp := time.Now().Unix()
	testConfig := fmt.Sprintf(`
	data "transip_domain" "test" {
		name = "%s"
	}

	resource "transip_dns_record" "test1" {
		domain  = "${data.transip_domain.test.id}"
		name    = "_terraform_provider_transip1_%d"
		type    = "CNAME"
		content = ["@"]
	}
	`, os.Getenv("TRANSIP_TEST_DOMAIN"), timestamp)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("transip_dns_record.test1", "name"),
				),
			},
		},
	})
}

func TestAccTransipResourceDomainConcurrent(t *testing.T) {
	timestamp := time.Now().Unix()
	testConfig := fmt.Sprintf(`
	data "transip_domain" "test" {
		name = "%s"
	}

	resource "transip_dns_record" "test1" {
		domain  = "${data.transip_domain.test.id}"
		name    = "_terraform_provider_transip1_%d"
		type    = "CNAME"
		content = ["@"]
	}

	resource "transip_dns_record" "test2" {
		domain  = "${data.transip_domain.test.id}"
		name    = "_terraform_provider_transip2_%d"
		type    = "CNAME"
		content = ["@"]
	}
	`, os.Getenv("TRANSIP_TEST_DOMAIN"), timestamp, timestamp)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("transip_dns_record.test1", "name"),
				),
			},
		},
	})
}
