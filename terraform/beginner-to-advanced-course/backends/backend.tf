terraform {
  backend "s3" {
    bucket         = "terraform-course-backend-s26"
    key            = "terraform.tfstate"
    region         = "eu-west-1"
    dynamodb_table = ""
  }
}