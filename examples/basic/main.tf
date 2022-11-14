terraform {
  required_providers {
    cohesivenet = {
      #source = "cohesive/cohesivenet"
      source = "cohesive.net/vns3/cohesivenet"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

# Configure the AWS Provider
provider "aws" {
  region = "eu-north-1"
}

/* Configure Cohesive Terraform Provider

We are creating a new controller in this script
so don't know controller address yet, leave empty.
Note: If this script was executing on a licensed controller,
then configure provider with address as follows:

provider "cohesivenet" {
  vns3 {
    username = "api"
    password = "password for api user on vns3 controller"
    host = "ip address of controller"
  }
}
*/
provider "cohesivenet" {

}

/******  Configure variables for VNS3 controller **********/

variable "vns3_subnet_id" {
    default = "subnet-081a2948e156a6c30"
}

variable "vns3_security_group_id" {
  default = "sg-055a5083372c44005"
}

variable "vns3_ami" {
    default = "ami-0b18e43060394578f"
}

variable "vns3_instance_type" {
  default = "t3.small"
}

variable "vns3_controller_name" {
    default = "terraform"
}

variable "vns3_master_password" {
  type = string
  default = "mypassword"
  sensitive =  true
}

variable "vns3_license_file" {
  default = "~/license.txt"
}

variable "keyset_token" {
    default = "keysetpassword"
    sensitive =  true
}

/****** End of configure variables for VNS3 controller **********/


/* Create an AWS network interface */
resource "aws_network_interface" "vns3_controller_eni_primary" {
  subnet_id         = var.vns3_subnet_id
  private_ips_count = 1
  source_dest_check = false
  security_groups   = [
      var.vns3_security_group_id
  ] 
  tags              =  {
                          Name = format("controller-eni-%s", var.vns3_controller_name)
                        }
}

/* Create an AWS EC2 instance for VNS3 Controller */
resource "aws_instance" "vns3_controller" {
  ami               = var.vns3_ami
  instance_type     = var.vns3_instance_type
  tags              =  {
                           Name = format("vns3-%s", var.vns3_controller_name)
                        }
                        
    network_interface {
        network_interface_id = aws_network_interface.vns3_controller_eni_primary.id
        device_index         = 0
    }

    depends_on = [
        aws_network_interface.vns3_controller_eni_primary
    ]

}

/* Create an AWS EIP for VNS3 controller */
resource "aws_eip" "vns3_ip" {
  vpc               = true
  instance          = aws_instance.vns3_controller.id
  network_interface = aws_network_interface.vns3_controller_eni_primary.id
}

/* Configure properties for VNS3 controller */
resource "cohesivenet_vns3_config" "vns3" {
  vns3 {
    host = aws_eip.vns3_ip.public_ip
    password = aws_instance.vns3_controller.id
  }

  topology_name = "${var.vns3_controller_name} topology"
  controller_name = var.vns3_controller_name
  license_file = var.vns3_license_file
  new_api_password = var.vns3_master_password
  new_ui_password = var.vns3_master_password
  generate_token = false

  license_params {
      default = true
  }
  keyset_params {
      token = var.keyset_token
  }
  peer_id = 1
}

/* Add a default route for VNS3 controller 1 */
resource "cohesivenet_vns3_routes" "route" {
    vns3 {
      username = "api"
      host = aws_eip.vns3_ip.public_ip
      password = var.vns3_master_password 
    } 

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
  vns3 {
      username = "api"
      host = aws_eip.vns3_ip.public_ip
      password = var.vns3_master_password 
  } 
  rule {
    script = "PREROUTING -d 10.18.0.65 -p udp --dport 162 -j DNAT --to 198.52.100.5:162"
  } 

  depends_on = [
      cohesivenet_vns3_config.vns3
  ]  
}


/* Configure an nginx image for VNS3 controller */
/*
resource "cohesivenet_vns3_plugin_images" "nginx" {
  image {
    image_name = "nginx"
    url  = "https://st-temp-vf-share.s3.eu-central-1.amazonaws.com/nginx_lb_release_1509.export.tar.gz" 
  }
}
*/

/* Configure an ha image for VNS3 controller */
/*
resource "cohesivenet_vns3_plugin_images" "ha" {
  image {
    image_name = "ha"
    url  = "https://cohesive-networks-plugins.s3.amazonaws.com/plugins/vns3-high-availability/vns3-high-availability.pm.v2.2.1.tgz"
  }
}
*/

/*
locals {
  all_images = {
      "nginx" = { image_name = "nginx", url = "https://st-temp-vf-share.s3.eu-central-1.amazonaws.com/nginx_lb_release_1509.export.tar.gz" },
      "ha" = { image_name = "ha", url = "https://cohesive-networks-plugins.s3.amazonaws.com/plugins/vns3-high-availability/vns3-high-availability.pm.v2.2.1.tgz" },
  }
}

resource "cohesivenet_vns3_plugin_images" "images" {
  for_each = local.all_images
  image {
    image_name = each.value.image_name
    url = each.value.url
  }
}*/






