variable "vns3_master_password" {
  sensitive = true
  type = string
  default = "testtest"
}

variable "vns3_license_file" {
    # ADD PATH TO YOUR LICENSE FILE
  default = "/Users/foo/license.txt"
}

variable "keyset_token" {
    default = "testtest"
    sensitive =  true
}

variable "vns3_instance_type" {
  default = "t3.small"  
}

variable subnet_id_1 {
  default = "subnet-asfsdafsadfadsf"
}

variable subnet_id_2 {
  default = "subnet-asdfdasfdsafsafsa"
}

variable "security_group_id" {
    default = "sg-asfdsafsadfsda"
}

variable "ami_id_1" {
    default = "ami-asdfdsafsdafafsad"
}

variable "ami_id_2" {
    default = "ami-asdfadsfsafsdaf"
}
