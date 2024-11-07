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
  default = "../sme6.txt" #Contact Cohesive Networks for a PoC license file
}

variable "keyset_token" {
  default   = "testkeysettoken"
  sensitive = true
}

variable "vns3_instance_type" {
  default = "t3.small"
}

variable "vpc_cidr" {
  default = "10.100.0.0/24"
}

variable "az1" {
  default = "us-east-1a"
}

variable "az2" {
  default = "us-east-1b"
}

variable "ami_id_1" {
  //set in tfvars
}

variable "ami_id_2" {
  //set in tfvars
}
