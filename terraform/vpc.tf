resource "aws_default_vpc" "default" {
  tags = {
    "Name" = "Default VPC"
  }
}

resource "aws_default_subnet" "default_az1" {
  availability_zone = "us-east-1a"
}

resource "aws_default_subnet" "default_az2" {
  availability_zone = "us-east-1b"
}

resource "aws_default_subnet" "default_az3" {
  availability_zone = "us-east-1c"
}