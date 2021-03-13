provider "transip" { }
terraform {
  required_providers {
    transip = {
      source  = "aequitas/transip"
     }
  }
}