resource "aws_eip" "myip" {
  vpc = "true"
}

resource "aws_eip" "myip01" {
  vpc      = "true"
  provider = aws.central
}