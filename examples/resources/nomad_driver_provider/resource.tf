terraform {
  required_providers {
    nomad-driver = {
      source  = "hashicorp/nomad-driver"
      version = "1.0.0"
    }
}
}
provider "nomad-driver" {
  address = "http://localhost:4646"
}
resource "nomad-driver" "dotnet_driver" {}