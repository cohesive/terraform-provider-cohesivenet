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

/* Configure Cohesive Terraform Provider */
provider "cohesivenet" {
  username = "api"
  password = var.vns3_master_password
  host = aws_eip.vns3_ip.public_ip
}

/******  Configure variables for VNS3 controller **********/

variable "vns3_controller_name" {
    default = "VNS3-Terraform"
}

variable "vns3_ami" {
  default = "ami-0234234534527234e"
}

variable "vns3_instance_type" {
  default = "t3.small"
}

variable "vns3_master_password" {
  default = "testing-1"
  sensitive =  true
}

variable "vns3_license_file" {
  default = "/Users/foo/license.txt"
}

variable "keyset_token" {
    default = "keysetpassword"
    sensitive =  true
}

/****** End of configure variables for VNS3 controller **********/

/* Create an AWS EC2 instance for VNS3 Controller */
resource "aws_instance" "vns3_controller" {
  ami               = "ami-aslkdfklasldfjsdla"
  instance_type     = "t3.small"
  tags              =  { Name = "VNS3 - Terraform" }
                        
  lifecycle {
    create_before_destroy = true
  }
}

/* Create an AWS EIP for VNS3 controller */
resource "aws_eip" "vns3_ip" {
  vpc               = true
  instance          = aws_instance.vns3_controller.id
}

/* Configure properties for VNS3 controller */
resource "cohesivenet_vns3_config" "vns3" {
  vns3 {
    host = aws_eip.vns3_ip.public_ip
    password = aws_instance.vns3_controller.id
  }
  configuration_id = aws_instance.vns3_controller.ami
  topology_name = "top-1"
  controller_name = "VNS3"
  license_file = "/Users/foo/license.txt"
  new_api_password = var.vns3_master_password
  new_ui_password = var.vns3_master_password
  license_params {
    default = true
  }
  keyset_params {
      token = "test-tf"
  }
  peer_id = 1
}

/* Add a default route for VNS3 controller 1 */
resource "cohesivenet_vns3_routes" "route" {
  route {
      cidr = "192.168.54.0/24"
      description = "default route"
      interface = ""
      gateway = "192.168.54.1/32"
      advertise = true
      metric = 300
  }
   
    depends_on = [
      cohesivenet_vns3_config.vns3
    ]  
}


resource "cohesivenet_vns3_firewall_rules" "rule" {
  rule {
    script = "PREROUTING -d 10.18.0.65 -p udp --dport 162 -j DNAT --to 198.52.100.5:162"
  } 

  depends_on = [
      cohesivenet_vns3_config.vns3
  ]  
}

