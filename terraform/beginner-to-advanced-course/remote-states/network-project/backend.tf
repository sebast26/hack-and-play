terraform {
  backend "s3" {
    bucket = "terraform-course-backend-s26"
    key    = "network/eip.tfstate"
    region = "eu-west-1"
  }
}