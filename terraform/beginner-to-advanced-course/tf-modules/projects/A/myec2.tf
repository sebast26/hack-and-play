module "ec2module" {
  source        = "../../modules/ec2"
  instance_type = "t3.micro"
}

resource "aws_instance" "myec2" {
  ami                    = "some-id"
  instance_type          = var.instance_type
  vpc_security_group_ids = [module.ec2module.sg_id]
}