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
    Name = "VNS3 C1"
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
    aws_eip.vns3_ip_1
  ]

}


