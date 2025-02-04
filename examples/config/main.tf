terraform {
  required_providers {
    cohesivenet = {
      version = "0.1.14"
      source = "cohesive/cohesivenet"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

# Configure the AWS Provider
provider "aws" {
  region = "us-east-1"
}

resource "aws_instance" "vns3controller_1" {
  ami               = var.ami_id_1
  instance_type     = "t3.small"
  subnet_id       = var.subnet_id_1
  source_dest_check = false
  tags              =  { Name = "Controller-1" }
  vpc_security_group_ids = [
      var.security_group_id
  ]
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_instance" "vns3controller_2" {
  ami               = var.ami_id_2
  instance_type     = "t3.small"
  subnet_id       = var.subnet_id_2
  source_dest_check = false
  tags              =  { Name = "Controller-2" }
  vpc_security_group_ids = [
      var.security_group_id
  ]
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_eip" "vns3_ip_1" {
  instance          = aws_instance.vns3controller_1.id
}

resource "aws_eip" "vns3_ip_2" {
  instance          = aws_instance.vns3controller_2.id
}

provider "cohesivenet" {
  alias    = "controller_1"
  host     = aws_eip.vns3_ip_1.public_ip
  password = aws_instance.vns3controller_1.id
}


provider "cohesivenet" {
  alias    = "controller_2"
  host     = aws_eip.vns3_ip_2.public_ip
  password = aws_instance.vns3controller_2.id
}

resource "cohesivenet_vns3_config" "vns3_1" {
  provider         = cohesivenet.controller_1
  configuration_id = var.ami_id_1
  instance_id  = aws_instance.vns3controller_1.id
  topology_name = "topology-1"
  controller_name = "controller-1"
  license_file = "/Users/foo/license.txt"
  new_api_password = var.vns3_master_password
  new_ui_password = var.vns3_master_password
  generate_token = true
  token_lifetime = 86400
  token_refresh = true
  license_params {
      default = true
  }
  keyset_params {
      token = var.keyset_token
      //no source needed as this will be licensed
  }
  peer_id = 1
  depends_on = [
    aws_instance.vns3controller_1,
    aws_eip.vns3_ip_1
  ]
}

resource "cohesivenet_vns3_config" "vns3_2" {
  provider         = cohesivenet.controller_2
  instance_id      = aws_instance.vns3controller_1.id
  configuration_id = var.ami_id_2
  topology_name = "topology-2"
  controller_name = "controller-2"
  new_api_password = var.vns3_master_password
  new_ui_password = var.vns3_master_password
  generate_token = true
  token_lifetime = 86400
  token_refresh = true
  license_params {
      default = true
  }
  keyset_params {
      token = var.keyset_token
      source = aws_instance.vns3controller_1.private_ip
  }
  peer_id = 2
  depends_on = [
    cohesivenet_vns3_config.vns3_1,
    aws_instance.vns3controller_2,
    aws_eip.vns3_ip_2
  ]
}

resource "cohesivenet_vns3_peers" "vns3_1_peers" {
  peer {
    address = aws_instance.vns3controller_2.private_ip
    peer_id = 2
  }
  depends_on = [
    cohesivenet_vns3_config.vns3_1,
    cohesivenet_vns3_config.vns3_2
  ]
}

resource "cohesivenet_vns3_peers" "vns3_2_peers" {
  peer {
    address = aws_instance.vns3controller_1.private_ip
    peer_id = 1
  }
  depends_on = [
    cohesivenet_vns3_config.vns3_1,
    cohesivenet_vns3_config.vns3_2
  ]
}