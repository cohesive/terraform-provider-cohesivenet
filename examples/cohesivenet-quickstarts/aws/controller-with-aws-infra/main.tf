terraform {
  required_providers {
    cohesivenet = {
      source  = "cohesive/cohesivenet"
      version = "1.0.5"
    }
    aws = {
      version = "0.1.0"
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

# Configure the AWS Provider
provider "aws" {
  region = "us-east-1"
}
resource "aws_vpc" "controller_tests" {
  cidr_block       = var.vpc_cidr
  instance_tenancy = "default"

  tags = {
    Name = "cohesive-examples"
  }
}

resource "aws_subnet" "public_1" {
  vpc_id            = aws_vpc.controller_tests.id
  cidr_block = "${cidrsubnet(aws_vpc.controller_tests.cidr_block, 4, 0)}"
  availability_zone = var.az1

  tags = {
    Name = "cohesive-examples-public-1"
  }
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.controller_tests.id

  tags = {
    Name = "cohesive-examples-igw"
  }
}

resource "aws_route_table" "controller_tests_route_table" {
  vpc_id = aws_vpc.controller_tests.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }

  tags = {
    Name = "cohesive-examples-rt"
  }
}

resource "aws_route_table_association" "public_1" {
  route_table_id = aws_route_table.controller_tests_route_table.id
  subnet_id      = aws_subnet.public_1.id

  depends_on = [aws_route_table.controller_tests_route_table]
}


resource "aws_security_group" "controller_tests_security_group" {
  name        = "cohesive-examples-sg"
  description = "cohesive-examples-sg"
  vpc_id      = aws_vpc.controller_tests.id

  ingress {
    description = "UI"
    from_port   = 8000
    to_port     = 8000
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "Peering"
    from_port   = 1194
    to_port     = 1203
    protocol    = "udp"
    cidr_blocks = [aws_vpc.controller_tests.cidr_block]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    Name = "cohesive-examples-sg"
  }
}

provider "cohesivenet" {
  alias    = "controller_1"
  host     = aws_eip.vns3_ip_1.public_ip
  password = aws_instance.vns3controller_1.id
}

resource "aws_instance" "vns3controller_1" {

  ami               = var.ami_id_1
  instance_type     = var.vns3_instance_type
  subnet_id         = aws_subnet.public_1.id
  source_dest_check = false
  vpc_security_group_ids = [
    aws_security_group.controller_tests_security_group.id
  ]
  tags = {
    Name = "Cohesive VNS3 C1"
  }
  lifecycle {
    create_before_destroy = true
  }

}

resource "aws_eip" "vns3_ip_1" {
  instance = aws_instance.vns3controller_1.id
}

resource "cohesivenet_vns3_config" "vns3_1" {
  provider         = cohesivenet.controller_1
  configuration_id = var.ami_id_1
  instance_id      = aws_instance.vns3controller_1.id
  topology_name    = var.topology_name
  controller_name  = var.controller_name
  license_file     = var.vns3_license_file
  #new_api_password = var.vns3_master_password
  #new_ui_password = var.vns3_master_password

  license_params {
    default = true
  }
  keyset_params {
    token = var.keyset_token
  }
  peer_id = 1

  depends_on = [
    aws_instance.vns3controller_1,
    aws_eip.vns3_ip_1,
    aws_route_table.controller_tests_route_table,
    aws_security_group.controller_tests_security_group,
    aws_route_table_association.public_1
  ]

}