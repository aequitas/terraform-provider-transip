package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTransipResourceDomain(t *testing.T) {
	if v := os.Getenv("TF_VAR_domain"); v == "" {
		t.Skip("TF_VAR_domain must be set for acceptance tests")
	}

	timestamp := time.Now().Unix()
	testConfig := fmt.Sprintf(`
	data "transip_domain" "test" {
		name = "%s"
	}

	resource "transip_dns_record" "test1" {
		domain  = data.transip_domain.test.id
		name    = "terraform-provider-transip1-%d"
		type    = "CNAME"
		content = ["@"]
	}
	`, os.Getenv("TF_VAR_domain"), timestamp)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("transip_dns_record.test1", "content.#", "1"),
				),
			},
		},
	})
}

func TestAccTransipResourceDomainMultiple(t *testing.T) {
	if v := os.Getenv("TF_VAR_domain"); v == "" {
		t.Skip("TF_VAR_domain must be set for acceptance tests")
	}

	timestamp := time.Now().Unix()
	testConfig := fmt.Sprintf(`
	data "transip_domain" "test" {
		name = "%s"
	}

	resource "transip_dns_record" "test2" {
		domain  = data.transip_domain.test.id
		name    = "terraform-provider-transip2-%d"
		type    = "A"
		content = ["192.0.2.0", "192.0.2.1", "192.0.2.2", "192.0.2.3"]
	}
	`, os.Getenv("TF_VAR_domain"), timestamp)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("transip_dns_record.test2", "content.#", "4"),
				),
			},
		},
	})
}

func TestAccTransipResourceDomainUpdate(t *testing.T) {
	if v := os.Getenv("TF_VAR_domain"); v == "" {
		t.Skip("TF_VAR_domain must be set for acceptance tests")
	}

	timestamp := time.Now().Unix()
	testConfig := fmt.Sprintf(`
	terraform { required_version = ">= 0.12.0" }

	data "transip_domain" "test" {
		name = "%s"
	}

	resource "transip_dns_record" "test7" {
		domain  = data.transip_domain.test.id
		name    = "terraform-provider-transip7-%d"
		type    = "A"
		content = ["192.0.2.0", "192.0.2.1"]
	}
  `, os.Getenv("TF_VAR_domain"), timestamp)
	testConfig2 := fmt.Sprintf(`
	terraform { required_version = ">= 0.12.0" }

	data "transip_domain" "test" {
		name = "%s"
	}

	resource "transip_dns_record" "test7" {
		domain  = data.transip_domain.test.id
		name    = "terraform-provider-transip7-%d"
		type    = "A"
		content =  ["192.0.2.2", "192.0.2.3"]
	}
	`, os.Getenv("TF_VAR_domain"), timestamp)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("transip_dns_record.test7", "content.#", "2"),
				),
			},
			{
				Config: testConfig2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("transip_dns_record.test7", "content.#", "2"),
				),
			},
		},
	})
}

func TestAccTransipResourceDomainConcurrent(t *testing.T) {
	if v := os.Getenv("TF_VAR_domain"); v == "" {
		t.Skip("TF_VAR_domain must be set for acceptance tests")
	}

	if v := os.Getenv("TF_VAR_domain"); v == "" {
		t.Skip("TF_VAR_domain must be set for acceptance tests")
	}

	timestamp := time.Now().Unix()
	testConfig := fmt.Sprintf(`
  terraform { required_version = ">= 0.12.0" }

  data "transip_domain" "test" {
    name = "%s"
  }

  resource "transip_dns_record" "test3" {
    domain  = data.transip_domain.test.id
    name    = "terraform-provider-transip3-%d"
    type    = "CNAME"
    content = ["@"]
  }

  resource "transip_dns_record" "test4" {
    domain  = data.transip_domain.test.id
    name    = "terraform-provider-transip4-%d"
    type    = "CNAME"
    content = ["@"]
  }
  `, os.Getenv("TF_VAR_domain"), timestamp, timestamp)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("transip_dns_record.test3", "content.#", "1"),
					resource.TestCheckResourceAttr("transip_dns_record.test4", "content.#", "1"),
				),
			},
		},
	})
}

func TestAccTransipResourceDomainConcurrentMultiple(t *testing.T) {
	domain := os.Getenv("TF_VAR_domain")
	if domain == "" {
		t.Skip("TF_VAR_domain must be set for acceptance tests")
	}

	timestamp := time.Now().Unix()
	testConfig := fmt.Sprintf(`
  terraform { required_version = ">= 0.12.0" }

  data "transip_domain" "test" {
    name = "%s"
  }

  resource "transip_dns_record" "test5" {
    domain  = data.transip_domain.test.id
    name    = "terraform-provider-transip5-%d"
    type    = "A"
    content = ["192.0.2.0", "192.0.2.1"]
  }

  resource "transip_dns_record" "test6" {
    domain  = data.transip_domain.test.id
    name    = "terraform-provider-transip6-%d"
    type    = "A"
    content =  ["192.0.2.2", "192.0.2.3"]
  }
  `, domain, timestamp, timestamp)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("transip_dns_record.test5", "content.#", "2"),
					resource.TestCheckResourceAttr("transip_dns_record.test6", "content.#", "2"),
				),
			},
		},
	})
}

func TestAccTransipResourceDomainRenameRecord(t *testing.T) {
	if v := os.Getenv("TF_VAR_domain"); v == "" {
		t.Skip("TF_VAR_domain must be set for acceptance tests")
	}

	timestamp := time.Now().Unix()
	testConfig := fmt.Sprintf(`
	terraform { required_version = ">= 0.12.0" }

	data "transip_domain" "test" {
		name = "%s"
	}

	resource "transip_dns_record" "test" {
		domain  = data.transip_domain.test.id
		name    = "terraform-provider-transip-%d"
		type    = "A"
		content = ["192.0.2.0", "192.0.2.1"]
	}
  `, os.Getenv("TF_VAR_domain"), timestamp)
	testConfig2 := fmt.Sprintf(`
	terraform { required_version = ">= 0.12.0" }

	data "transip_domain" "test" {
		name = "%s"
	}

	resource "transip_dns_record" "test" {
		domain  = data.transip_domain.test.id
		name    = "terraform-provider-transip-changed-%d"
		type    = "A"
		content = ["192.0.2.0", "192.0.2.1"]
	}
	`, os.Getenv("TF_VAR_domain"), timestamp)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("transip_dns_record.test", "content.#", "2"),
				),
			},
			{
				Config: testConfig2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("transip_dns_record.test", "content.#", "2"),
				),
			},
		},
	})
}

func TestAccTransipResourceDomainChangeType(t *testing.T) {
	if v := os.Getenv("TF_VAR_domain"); v == "" {
		t.Skip("TF_VAR_domain must be set for acceptance tests")
	}

	timestamp := time.Now().Unix()
	testConfig := fmt.Sprintf(`
	terraform { required_version = ">= 0.12.0" }

	data "transip_domain" "test" {
		name = "%s"
	}

	resource "transip_dns_record" "test" {
		domain  = data.transip_domain.test.id
		name    = "terraform-provider-transip-%d"
		type    = "A"
		content = ["192.0.2.0", "192.0.2.1"]
	}
  `, os.Getenv("TF_VAR_domain"), timestamp)
	testConfig2 := fmt.Sprintf(`
	terraform { required_version = ">= 0.12.0" }

	data "transip_domain" "test" {
		name = "%s"
	}

	resource "transip_dns_record" "test" {
		domain  = data.transip_domain.test.id
		name    = "terraform-provider-transip-changed-%d"
		type    = "CNAME"
		content = ["example.com."]
	}
	`, os.Getenv("TF_VAR_domain"), timestamp)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("transip_dns_record.test", "content.#", "2"),
				),
			},
			{
				Config: testConfig2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("transip_dns_record.test", "content.#", "1"),
				),
			},
		},
	})
}
