# Domain Nameservers Resource

Allows to set non default (ie. non TransIP provide) nameservers.

Use together with `domain_dnssec` to safely use another DNS Provider (ie
Cloudflare, AWS Route53, ....)

## Argument Reference

* `domain` - (Required) The domain, including the tld
* `nameserver` - (Required) List of nameservers associated with domain

### Nameserver object

* `hostname` - (Required) The hostname of this nameserver
* `ipv4` - (Optional) ipv4 glue record for this nameserver
* `ipv6` - (Optional) ipv6 glue record for this nameserver

## Attribute Reference

* `id` - n/a
