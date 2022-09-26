terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = ">=3.60,<=3.70"
    }
  }
}

provider "aws" {
  # Configuration options
}