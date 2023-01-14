terraform {
  required_version = ">= 1.0.0"
}

resource "digitalocean_droplet" "w10k-go" {
  image      = "ubuntu-22-10-x64"
  name       = "w10k-go"
  region     = "nyc1"
  size       = "s-1vcpu-1gb"
  ssh_keys   = [data.digitalocean_ssh_key.do.id]
  monitoring = true

  connection {
    host        = self.ipv4_address
    user        = "root"
    type        = "ssh"
    private_key = file(var.pvt_key)
    timeout     = "2m"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo apt-get update",
    ]
  }
}

resource "digitalocean_domain" "default" {
  name       = format("w10k-go.%s", var.domain)
  ip_address = digitalocean_droplet.w10k-go.ipv4_address
}

