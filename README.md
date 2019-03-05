# Terraform Transip provider

Provides resources for Transip resources using [Transip API](https://www.transip.nl/transip/api/)

Supported resources:

    - Domain name (data source, resource)
    - Domain name DNS records (resource)

## Requirements

In order to use the provider you need a Transip account. For this account the API should be enabled and a private key should be created which is used for authentication (https://www.transip.nl/cp/account/api/).

## Installation

Download the latest binary release from the [Releases](https://github.com/aequitas/terraform-provider-transip/releases) page, unzip it to a location in `PATH` (eg: `/usr/local/bin/`).

## Example

```hcl
variable "private_key" {}

provider "transip" {
  account_name = "example"
  private_key  = "${var.private_key}"
}

data "transip_domain" "example_com" {
  name = "example.com"
}

resource "transip_dns_record" "www" {
  domain  = "${transip_domain.example_com.id}"
  name    = "www"
  type    = "CNAME"
  content = ["@"]
}

resource "transip_dns_record" "test" {
  domain = "${transip_domain.example_com.id}"
  name   = "test"
  type   = "A"

  content = [
    "203.0.113.1",
    "203.0.113.2",
  ]
}

resource "transip_dns_record" "test" {
  domain = "${transip_domain.example_com.id}"
  name   = "test"
  expire = 300
  type   = "AAAA"

  content = [
    "2001:db8::1",
  ]
}
```

## Roadmap

- Tests
- VPS (data source, resource)
