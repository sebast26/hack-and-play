variable "environment" {
  type        = string
  description = "Infrastructure environment. eg. dev, prod, etc"
  default     = "test"
}

variable "vpc_name" {
  type    = string
  default = "test_vpc"
}