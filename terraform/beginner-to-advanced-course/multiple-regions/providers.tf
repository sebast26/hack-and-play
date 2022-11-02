provider "aws" {
  region = "eu-west-1"
}

provider "aws" {
  alias  = "central"
  region = "eu-central-1"
  # if you have other AWS account
  # profile = "second aws profile"
}