provider "aws" {
  region = "eu-west-1"
}

module "ec2-instance" {
  source  = "terraform-aws-modules/ec2-instance/aws"
  version = "4.1.4"

  // variables that are required
  // see documentation
}