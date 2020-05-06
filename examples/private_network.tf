# order a new private network, or use `terraform import transip_private_netwerk.test test` to import existing one
resource "transip_private_network" "test" {
  description = "test" # Description is the user defined name, as the name attribute is set by TransIP as a unique ID.
}