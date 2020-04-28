# Because concurrency is a bit finicky at the moment run this with -parallelism=1

variable "domain" {
  default = "example.com"
}

data "transip_domain" "demo" {
  name = var.domain
}

resource "transip_dns_record" "www" {
  domain  = data.transip_domain.demo.id
  name    = "www"
  type    = "CNAME"
  content = ["@"]
}

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

resource "transip_domain" "test" {
  name = var.domain
}