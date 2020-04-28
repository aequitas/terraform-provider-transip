package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"os"
	"regexp"
	"testing"
	"time"
)

func TestAccTransipResourceDomain(t *testing.T) {
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
					resource.TestCheckResourceAttr("transip_dns_record.test1", "content.0", "@"),
				),
			},
		},
	})
}

func TestAccTransipResourceDomainMultiple(t *testing.T) {
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
					resource.TestMatchResourceAttr("transip_dns_record.test2", "content.0", regexp.MustCompile("192.0.2.[0-3]")),
				),
			},
		},
	})
}

func TestAccTransipResourceDomainUpdate(t *testing.T) {
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
					resource.TestMatchResourceAttr("transip_dns_record.test7", "content.0", regexp.MustCompile("192.0.2.[01]")),
					resource.TestMatchResourceAttr("transip_dns_record.test7", "content.1", regexp.MustCompile("192.0.2.[01]")),
				),
			},
			{
				Config: testConfig2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("transip_dns_record.test7", "content.0", regexp.MustCompile("192.0.2.[23]")),
					resource.TestMatchResourceAttr("transip_dns_record.test7", "content.1", regexp.MustCompile("192.0.2.[23]")),
				),
			},
		},
	})
}

// // TODO: concurrency seems broken on the Transip side, needs more testing to prove
// func TestAccTransipResourceDomainConcurrent(t *testing.T) {
//   timestamp := time.Now().Unix()
//   testConfig := fmt.Sprintf(`
//   terraform { required_version = ">= 0.12.0" }
//
//   data "transip_domain" "test" {
//     name = "%s"
//   }
//
//   resource "transip_dns_record" "test3" {
//     domain  = data.transip_domain.test.id
//     name    = "terraform-provider-transip3-%d"
//     type    = "CNAME"
//     content = ["@"]
//   }
//
//   resource "transip_dns_record" "test4" {
//     domain  = data.transip_domain.test.id
//     name    = "terraform-provider-transip4-%d"
//     type    = "CNAME"
//     content = ["@"]
//   }
//   `, os.Getenv("TF_VAR_domain"), timestamp, timestamp)
//
//   resource.Test(t, resource.TestCase{
//     PreCheck:  func() { testAccPreCheck(t) },
//     Providers: testAccProviders,
//     Steps: []resource.TestStep{
//       {
//         Config: testConfig,
//         Check: resource.ComposeTestCheckFunc(
//           resource.TestCheckResourceAttr("transip_dns_record.test3", "content.0", "@"),
//           resource.TestCheckResourceAttr("transip_dns_record.test4", "content.0", "@"),
//         ),
//       },
//     },
//   })
// }
//
// // TODO: concurrency seems broken on the Transip side, needs more testing to prove
// func TestAccTransipResourceDomainConcurrentMultiple(t *testing.T) {
//   timestamp := time.Now().Unix()
//   testConfig := fmt.Sprintf(`
//   terraform { required_version = ">= 0.12.0" }
//
//   data "transip_domain" "test" {
//     name = "%s"
//   }
//
//   resource "transip_dns_record" "test5" {
//     domain  = data.transip_domain.test.id
//     name    = "terraform-provider-transip5-%d"
//     type    = "A"
//     content = ["192.0.2.0", "192.0.2.1"]
//   }
//
//   resource "transip_dns_record" "test6" {
//     domain  = data.transip_domain.test.id
//     name    = "terraform-provider-transip6-%d"
//     type    = "A"
//     content =  ["192.0.2.2", "192.0.2.3"]
//   }
//   `, os.Getenv("TF_VAR_domain"), timestamp, timestamp)
//
//   resource.Test(t, resource.TestCase{
//     PreCheck:  func() { testAccPreCheck(t) },
//     Providers: testAccProviders,
//     Steps: []resource.TestStep{
//       {
//         Config: testConfig,
//         Check: resource.ComposeTestCheckFunc(
//           resource.TestMatchResourceAttr("transip_dns_record.test5", "content.0", regexp.MustCompile("192.0.2.[01]")),
//           resource.TestMatchResourceAttr("transip_dns_record.test6", "content.0", regexp.MustCompile("192.0.2.[23]")),
//         ),
//       },
//     },
//   })
// }
