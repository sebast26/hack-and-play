provider "aws" {
  region = "eu-west-1"

  assume_role {
    role_arn     = ""
    session_name = "whatever"
  }
}
