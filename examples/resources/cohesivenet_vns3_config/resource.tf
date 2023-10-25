terraform {
  required_providers {
    cohesivenet = {
      version = "0.1.0"
      source  = "cohesive.net/vns3/cohesivenet"
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


provider "cohesivenet" {
  alias = "controller_1"
  username = var.username
  host = aws_eip.vns3_ip_1.public_ip
  password = aws_instance.vns3controller_1.id
  //token = cohesivenet_vns3_config.vns3_1.token
}


provider "cohesivenet" {
  alias = "controller_2"
  username = var.username
  host = aws_eip.vns3_ip_2.public_ip
  password = aws_instance.vns3controller_2.id
  //token = cohesivenet_vns3_config.vns3_2.token
}

resource "aws_instance" "vns3controller_1" {

  ami               = "var.ami_id_1"
  instance_type     = "var.instance_type"
  subnet_id       = "var.subnet_id"
  source_dest_check = false
  vpc_security_group_ids = [
      "var.security_group_id"
  ]
  tags = {
      Name = "var.controller_2_name"
    }
  lifecycle {
    create_before_destroy = true
  }

}

resource "aws_instance" "vns3controller_2" {

  ami               = "var.ami_id_2"
  instance_type     = "var.instance_type"
  subnet_id       = "var.subnet_id"
  source_dest_check = false
  vpc_security_group_ids = [
      "var.security_group_id"
  ]
  tags = {
    Name = "var.controller_1_name"
    }
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

resource "cohesivenet_vns3_config" "vns3_1" {
  provider = cohesivenet.controller_1
  configuration_id = var.ami_id_1
  instance_id = aws_instance.vns3controller_1.id
  topology_name = var.topology_name
  controller_name = var.controller_name
  license_file = var.license_file
  new_api_password = var.vns3_master_password
  new_ui_password = var.vns3_master_password
  generate_token = var.vns3_api_token_lifetime == 0 ? false : true
  token_lifetime = var.vns3_api_token_lifetime
  token_refresh = var.vns3_api_token_refresh
 

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


resource "cohesivenet_vns3_config" "vns3_2" {
  provider = cohesivenet.controller_2
  configuration_id = var.ami_id_2
  instance_id = aws_instance.vns3controller_2.id
  topology_name = var.topology_name
  controller_name = var.controller_name
  new_api_password = var.vns3_master_password
  new_ui_password = var.vns3_master_password
  generate_token = var.vns3_api_token_lifetime == 0 ? false : true
  token_lifetime = var.vns3_api_token_lifetime
  token_refresh = var.vns3_api_token_refresh

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
provider = cohesivenet.controller_1
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
provider = cohesivenet.controller_2 
  peer {
    address = aws_instance.vns3controller_1.private_ip
    peer_id = 1
  }


  depends_on = [
    cohesivenet_vns3_config.vns3_1,
    cohesivenet_vns3_config.vns3_2
  ]
}