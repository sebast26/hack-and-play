provider "aws" {
    region = "eu-west-1"
}

locals {
    time = formatdate("DD MMM YYYY hh:mm ZZZ", timestamp())
}

variable "region" {
    default = "eu-west-1"
}

variable "tags" {
    type = list
    default = ["firstec2", "secondec2"]
}

variable "ami" {
    type = map
    default = {
        "eu-west-1" = "ami-0d71ea30463e0ff8d"
    }
}

resource "aws_key_pair" "loginkey" {
    key_name = "login-key"
    public_key = file("${path.module}/id_rsa.pub")
}

resource "aws_instance" "ec2-dev" {
    count = 2
    ami = lookup(var.ami, var.region, "no ami for region")
    instance_type = "t2.micro"
    key_name = aws_key_pair.loginkey.key_name

    tags = {
        Name = element(var.tags, count.index)
    }
}

output "timestamp" {
    value = local.time
}