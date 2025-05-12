# Terraform Transip provider

[![Build Status](https://travis-ci.org/aequitas/terraform-provider-transip.svg?branch=master)](https://travis-ci.org/aequitas/terraform-provider-transip)

Provides resources for Transip resources using [Transip API](https://www.transip.eu/transip/api/)

 - [Terraform Registry](https://registry.terraform.io/providers/aequitas/transip/latest)
 - [Terraform Registry Documentation](https://registry.terraform.io/providers/aequitas/transip/latest/docs)

Supported resources:

 - Domain name (data source, resource)
 - Domain name DNS records (resource)
 - VPS (data source, resource)
 - VPS private network (data source, resource)
 - VPS firewall (resource)
 - OpenStack Project (data source, resource)
 - OpenStack User (data source, resource)

## Requirements

In order to use the provider you need a Transip account. For this account the API should be enabled and a private key should be created which is used for authentication (https://www.transip.eu/cp/account/api/).

## Installation

### Terraform 0.13

Add the following snippet to a `.tf` somewhere in your codebase and run `terraform init` to automatically download and install the provider:

```hcl
terraform {
  required_providers {
    transip = {
      source = "aequitas/transip"
    }
  }
}
```

### Terraform <=0.12

Download the latest binary release from the [Releases](https://github.com/aequitas/terraform-provider-transip/releases) page, unzip it to one of these directories: [Plugin Locations](https://www.terraform.io/docs/extend/how-terraform-works.html#plugin-locations). Or refer to [Installing Plugins](https://www.terraform.io/docs/plugins/basics.html#installing-plugins) for more information.

## Notes

- The Transip API managed DNS Entries as a list property of a Domain object. In this implementation I have opted to give DNS entries their own resource `transip_dns_record` to make management more in line with other Terraform DNS Providers.

- Not all resources (especially the VPS resource) have been thoroughly tested. Use with care.

## Example

Also see examples in: [examples/](https://github.com/aequitas/terraform-provider-transip/tree/master/examples).
# Have Terraform install provider on init (Terraform 0.13 only)
```hcl
terraform {
  required_providers {
    transip = {
      source = "aequitas/transip"
    }
  }
}

# Enable Transip API, whitelist your IP, create private key and provide it here
provider "transip" {
  account_name = "example"
  private_key  = <<EOF
  -----BEGIN PRIVATE KEY-----
  ...
  -----END PRIVATE KEY-----
  EOF
}

# Or simply leave the provider empty when using the environment variables TRANSIP_ACCOUNT_NAME and TRANSIP_PRIVATE_KEY
# provider "transip" { }

# Get an existing domain as data source
data "transip_domain" "example_com" {
  name = "example.com"
}

# Or create/import a (new) domain name to be managed by Terraform
# resource "transip_domain" "example_com" {
#     name = "example.com"
# }

# Simple CNAME record
resource "transip_dns_record" "www" {
  domain  = data.transip_domain.example_com.id
  name    = "www"
  type    = "CNAME"
  content = ["@"]
}

# VPS Server with setup script and DNS record
resource "transip_vps" "test" {
  description = "example"
  product_name = "vps-bladevps-x2"
  operating_system = "ubuntu-18.04"

  # Script to run to provision the VPS
  install_text = <<EOF
  # install and enable firewall and basic webserver
  apt update
  apt install -yqq ufw nginx
  ufw allow 22/tcp
  ufw allow 80/tcp
  ufw allow 443/tcp
  ufw --force enable
  EOF
}
resource "transip_dns_record" "vps" {
  domain = data.transip_domain.example_com.id
  name   = "vps"
  type   = "A"

  content = [ transip_vps.test.ip_address ]
}

# A record with multiple entries, eg: for round robin DNS
resource "transip_dns_record" "test" {
  domain = data.transip_domain.example_com.id
  name   = "test"
  type   = "A"

  content = [
    "203.0.113.1",
    "203.0.113.2",
  ]
}

# IPv6 record
resource "transip_dns_record" "testv6" {
  domain = data.transip_domain.example_com.id
  name   = "test"
  expire = 300
  type   = "AAAA"

  content = [
    "2001:db8::1",
  ]
}

# Get an existing VPS as datasource
data "transip_vps" "test" {
  description = "example"
}

# Set hostname for VPS using data source
resource "transip_dns_record" "vps" {
  domain = data.transip_domain.example_com.id
  name   = "vps"
  type   = "A"

  content = [data.transip_vps.test.ip_address]
}
resource "transip_dns_record" "vps" {
  domain = data.transip_domain.example_com.id
  name   = "vps"
  type   = "AAAA"

  content = [data.transip_vps.test.ipv6_addresses[0]]
}
```

## Development

This project can be build and tested like any regular Go project or Terraform provider. For convenience a Makefile is provided which contains commands to easy recurring development tasks.

[Direnv](https://direnv.net/) and [keyring](https://pypi.org/project/keyring/) are used to setup environment variables used during testing and to keep credentials out of the project directory.

### Makefile

To test build the project simple run:

    make

As will setup dependencies, build binaries in `./build/` and install them to `./terraform.d/plugins/` so they can be used for testing.

When source files change te dependencies and binaries will automatically be rebuild (if required).

### Test suites

To run just the unit test suite (not requiring credentials or touching the API) run:

    make test

To run the acceptance test suite as well run:

    make test_acc

Note: the acceptance tests require Transip account credentials (username + certificate). The demo account won't work for this.

To configure these refer to `.envrc.local.example` file.

Warning: although care has been taken to prevent accidental modification of existing resource or unexpected costs to be made (by ordering product) this is not guaranteed. Use at own risk.

### Testing .tf files

To `plan` or `apply` the `.tf` files in `./examples/` you can run the following command:

    make plan

Or to apply:

    make apply

Source code is rebuild (if needed) and the plugin updated before running the Terraform command, allowing for quick iteration and debugging.

To just target a single resource (and it dependants) use the `targets` argument like so:

    make plan targets=transip_vps.test
