# Transip Provider

Provides resources for Transip resources using [Transip API](https://www.transip.eu/transip/api/)

## Example Usage

```hcl
# Have Terraform install provider on init (Terraform 0.13 only)
terraform {
  required_providers {
    transip = {
      source = "aequitas/transip"
    }
  }
}

provider "transip" {
  account_name = "example"
  private_key  = <<EOF
  -----BEGIN PRIVATE KEY-----
  ...
  -----END PRIVATE KEY-----
  EOF
}

# VPS Server with setup script and DNS record
resource "transip_vps" "test" {
  name = "example"
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

data "transip_domain" "example_com" {
  name = "example.com"
}

resource "transip_dns_record" "vps" {
  domain = data.transip_domain.example_com.id
  name   = "vps"
  type   = "A"

  content = [ transip_vps.test.ip_address ]
}
```

## Argument Reference

* `access_token` - (Optional) Temporary access token used for authentication.
* `account_name` - (Optional) Name of the Transip account.
* `private_key` - (Optional) Contents of the private key file to be used to authenticate.
* `read_only` - (Optional) Disable API write calls.
* `test_mode` - (Optional) Use API test mode.