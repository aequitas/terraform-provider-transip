variable "vps_name" {
  default = "test"
}

variable "vps_product" {
  default = "vps-bladevps-x2"
}

variable "vps_os" {
  default = "ubuntu-20.04"
}

# order a new VPS, or use `terraform import transip_vps.test test` to import existing one
resource "transip_vps" "test" {
  description = var.vps_name
  product_name = var.vps_product
  operating_system = var.vps_os
}