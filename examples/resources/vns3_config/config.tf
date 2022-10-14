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

provider "cohesivenet" {}

locals {
    vns3_version_parts = split("-", var.vns3_version)
    vns3_version_cln = replace(element(local.vns3_version_parts, 0), ".", "")
    vns3_version_date_regex = length(local.vns3_version_parts) > 1 ? element(local.vns3_version_parts, 1) : "[0-9a-z]+"
}

data "aws_ami" "vnscubed" {
    most_recent = true
    owners = ["${var.vns3_account_owner}"]
    name_regex = "^vnscubed${local.vns3_version_cln}-${local.vns3_version_date_regex}-${var.vns3_license_type}.*"

    filter {
        name = "root-device-type"
        values = ["ebs"]
    }
}

resource "aws_network_interface" "vns3controller_eni_primary" {
  count             = length(var.subnet_ids)
  subnet_id         = element(var.subnet_ids, count.index)
  private_ips_count = 1
  source_dest_check = false
  security_groups   = [
      var.security_group_id
  ] 
  tags              = merge(
                        var.common_tags,
                        {
                          Name = format("%s-controller-eni-%d", var.topology_name, count.index)
                        }
                      )
}

resource "aws_instance" "vns3controller" {
  ami               = data.aws_ami.vnscubed.id
  count             = length(var.subnet_ids)
  instance_type     = var.vns3_instance_type
  tags              = merge(
                        var.common_tags,
                        {
                          Name = format("%s-vns3-%d", var.topology_name, count.index)
                        }
                      )

    network_interface {
        network_interface_id = "${element(aws_network_interface.vns3controller_eni_primary.*.id, count.index)}"
        device_index         = 0
    }

    depends_on = [
        aws_network_interface.vns3controller_eni_primary
    ]

}

resource "aws_eip" "vns3_ip" {
  vpc               = true
  count             = length(aws_instance.vns3controller)
  instance          = element(aws_instance.vns3controller.*.id, count.index)
  network_interface = element(aws_network_interface.vns3controller_eni_primary.*.id, count.index)
}

resource "cohesivenet_vns3_config" "vns3_1" {
  vns3 {
    host = aws_eip.vns3_ip[0].public_ip
    password = aws_instance.vns3controller[0].id
  }

  topology_name = var.topology_name
  controller_name = "${var.controller_name} 1"
  license_file = var.vns3_license_file
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
}

resource "cohesivenet_vns3_config" "vns3_2" {
  vns3 {
    host = aws_eip.vns3_ip[1].public_ip
    password = aws_instance.vns3controller[1].id
  }

  topology_name = var.topology_name
  controller_name = "${var.controller_name} 2"
  new_api_password = var.vns3_master_password
  new_ui_password = var.vns3_master_password
  generate_token = var.vns3_api_token_lifetime == 0 ? false : true
  token_lifetime = var.vns3_api_token_lifetime
  token_refresh = var.vns3_api_token_refresh

  keyset_params {
      token = var.keyset_token
      source = aws_eip.vns3_ip[0].private_dns
  }

  peer_id = 2

  depends_on = [
    cohesivenet_vns3_config.vns3_1
  ]
}