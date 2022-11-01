data "terraform_remote_state" "eip" {
  backend = "s3"

  config = {
    bucket = "terraform-course-backend-s26"
    key    = "network/eip.tfstate"
    region = "eu-west-1"
  }
}