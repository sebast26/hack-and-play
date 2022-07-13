provider "aws" {
    region = "eu-west-1"
}

resource "aws_instance" "myec2" {
    ami = "ami-0d71ea30463e0ff8d"
    instance_type = "t2.micro"
}

provider "azurerm" {}

terraform {
  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
      version = "2.21.0"
    }
  }
}

provider "digitalocean" {}