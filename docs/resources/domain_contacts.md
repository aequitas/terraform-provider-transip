# Domain Contacts Resource

Manages whois contacts for a domain.

## Example Usage

```hcl
resource "transip_domain_contacts" "example" {
  domain = "example.com"

  contact {
    type         = "registrant"
    first_name   = "John"
    last_name    = "Doe"
    company_name = "Example Corp"
    street       = "Voorbeeldstraat"
    number       = "1"
    postal_code  = "1000 AA"
    city         = "Amsterdam"
    phone_number = "+31.201234567"
    email        = "john@example.com"
    country      = "nl"
  }
}
```

## Argument Reference

* `domain` - (Required) The domain name, including the tld.
* `contact` - (Required) List of whois contacts. At least one contact is required.

The `contact` block supports:

* `type` - (Required) The type of contact: `registrant`, `administrative` or `technical`.
* `first_name` - (Required) The first name.
* `last_name` - (Required) The last name.
* `company_name` - (Optional) The company name.
* `company_kvk` - (Optional) The KVK number.
* `company_type` - (Optional) The company type.
* `street` - (Required) The street name.
* `number` - (Required) The street number.
* `postal_code` - (Required) The postal code.
* `city` - (Required) The city.
* `phone_number` - (Required) The phone number.
* `fax_number` - (Optional) The fax number.
* `email` - (Required) The email address.
* `country` - (Required) The ISO 3166-1 alpha-2 country code (lowercase).

## Attribute Reference

* `id` - The domain name.
