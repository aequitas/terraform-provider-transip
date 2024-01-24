# Domain DNSSec Resource

Allows to set non default (ie. non TransIP provide) DNS Sec records.
Use together with `domain_nameservers` to safely use another DNS Provider (ie
Cloudflare, AWS Route53, ....)

## Argument Reference

* `domain` - (Required) The domain, including the tld
* `dnssec` - (Required) List of dnssec entries associated with domain

### DNSSec object

* `key_tag` - (Required) A 5-digit key of the Zonesigner
* `flags` - (Required) The signing key number, either 256 (Zone Signing Key) or 257 (Key Signing Key)
* `algorithm` - (Required) The algorithm type that is used, see: https://www.transip.nl/vragen/461-domeinnaam-nameservers-gebruikt-beveiligen-dnssec/ for the possible options.
* `public_key` - (Required) The public key

## Attribute Reference

* `id` - n/a
