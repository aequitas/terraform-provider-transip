# Because concurrency is a bit finicky at the moment run this with -parallelism=1

variable "domain" {
  default = "example.com"
}

# order domain, or use `terraform import transip_domain.test example.com` to import existing one
resource "transip_domain" "test" {
  name = var.domain
}

# use existing domain as data resource instead
data "transip_domain" "demo" {
  name = var.domain
}

# create a www CNAME record
resource "transip_dns_record" "www" {
  domain  = data.transip_domain.demo.id
  name    = "www"
  type    = "CNAME"
  content = ["@"]
}

# create a record with multiple addresses
resource "transip_dns_record" "demo" {
  domain  = data.transip_domain.demo.id
  name    = "demo"
  type    = "A"
  content = [
    "192.0.2.0",
    "192.0.2.1",
    "192.0.2.2",
    "192.0.2.3",
    "192.0.2.4",
    "192.0.2.5",
  ]
}

# use address of existing VPS as content for an A record
resource "transip_dns_record" "vps" {
  domain  = data.transip_domain.demo.id
  name    = "vps"
  type    = "A"
  content = [
    transip_vps.test.ip_address
  ]
}

# same but for IPv6
resource "transip_dns_record" "vps_v6" {
  domain  = data.transip_domain.demo.id
  name    = "vps"
  type    = "AAAA"
  content = [
    transip_vps.test.ipv6_addresses[0]
  ]
}