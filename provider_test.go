package main

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"transip": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("TRANSIP_ACCOUNT_NAME"); v == "" {
		t.Fatal("TRANSIP_ACCOUNT_NAME must be set for acceptance tests")
	}

	if v := os.Getenv("TRANSIP_PRIVATE_KEY"); v == "" {
		t.Fatal("TRANSIP_PRIVATE_KEY must be set for acceptance tests")
	}

	if v := os.Getenv("TF_VAR_domain"); v == "" {
		t.Fatal("TF_VAR_domain must be set for acceptance tests")
	}
}
