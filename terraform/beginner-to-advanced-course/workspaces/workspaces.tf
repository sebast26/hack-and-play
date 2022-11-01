resource "aws_instance" "myec2" {
  ami           = "ami-07e1daca3ee9095b3"
  instance_type = lookup(var.instance_type, terraform.workspace, "t2.mini")
}

variable "instance_type" {
  type = map(any)

  default = {
    default = "t2.nano"
    dev     = "t2.micro"
    prd     = "t2.large"
  }
}