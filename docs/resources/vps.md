# Vps Resource



## Argument Reference

* `availability_zone` - (Optional) The name of the availability zone the VPS is in.
* `description` - (Optional) The name that can be set by customer.
* `install_text` - (Optional) Base64 encoded preseed / kickstart / cloudinit instructions, when installing unattended. Provide plain text, not base64 encoded — the provider encodes it automatically.
* `install_flavour` - (Optional) The flavour of OS installation. Must be one of: `installer`, `preinstallable` or `cloudinit`. Defaults to `installer` when not specified. Use `cloudinit` when using cloud-init in `install_text`.
* `operating_system` - (Required) The VPS OperatingSystem.
* `product_name` - (Required) The product name.

## Attribute Reference

* `cpus` - The VPS cpu count.
* `disk_size` - The VPS disk size in kB.
* `id` - n/a
* `ip_address` - The VPS main ipAddress.
* `ipv4_addresses` - All IPV4 addresses associated with this VPS.
* `ipv6_addresses` - All IPV6 addresses associated with this VPS.
* `is_blocked` - If the VPS is administratively blocked.
* `is_customer_locked` - If this VPS is locked by the customer.
* `mac_address` - The VPS macaddress.
* `memory_size` - The VPS memory size in kB.
* `name` - The unique VPS name.
* `status` - The VPS status, either 'created', 'installing', 'running', 'stopped' or 'paused'.
* `tags` - The custom tags added to this VPS.