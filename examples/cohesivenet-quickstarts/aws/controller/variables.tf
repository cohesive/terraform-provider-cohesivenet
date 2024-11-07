/*
variable "vns3_master_password" {
  sensitive = true
  type      = string
}
*/
variable "topology_name" {
  default = "Tests"
}
variable "controller_name" {
  default = "VNS3"
}

variable "vns3_license_file" {
  default = "./vns3_free_license.txt"
}

variable "keyset_token" {
  default   = "testkeysettoken"
  sensitive = true
}

variable "vns3_instance_type" {
  default = "t3a.micro"
}

variable "vpc_cidr" {
  default = "10.100.0.0/24"
}

variable "az1" {
  default = "us-east-1a"
}


variable "ami_id_1" {
  //set in tfvars
}
