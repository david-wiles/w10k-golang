terraform {
  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

variable "do_token" {
  sensitive = true
}
variable "pvt_key" {
  sensitive = true
}
variable "domain" {}

provider "digitalocean" {
  token = var.do_token
}

data "digitalocean_ssh_key" "do" {
  name = "do.pub"
}
