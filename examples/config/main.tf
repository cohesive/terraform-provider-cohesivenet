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
  instance          = aws_instance.vns3controller[0].id
  network_interface = aws_network_interface.vns3controller_eni_primary[0].id
}

resource "cohesivenet_vns3_config" "vns3" {
    host = aws_eip.vns3_ip.public_ip
    password = aws_instance.vns3controller[0].id
    topology_name = var.topology_name
    controller_name = var.controller_name
    license_file = var.vns3_license_file
    license_params {
        default = true
    }
    keyset_params {
        token = var.keyset_token
    }
    peer_id = 1
}


#  resource "cohesivenet_routes" "route" {
#   route {
#     cidr = "192.168.54.0/24"
#     description = "cohesive_to_watford_secondary"
#     interface = "tun0"
#     gateway = "192.168.54.1/32"
#     advertise = true
#     metric = 300
#   }
#  }