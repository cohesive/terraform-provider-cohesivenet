variable "user" {
  type = string
}

variable "topology_name" {
    default = "bens-test-tf-launch"
}

variable "vns3_account_owner" {
  type = string
  default = "678554804139"
}

variable "vns3_version" {
    type = string
    default = "5.2.4-20220317"
}

variable "vns3_license_file" {
    default = "/Users/benplatta/code/cohesive/vns3-functional-testing/test-assets/license.txt"
}

variable "keyset_token" {
    default = "testtest"
    sensitive =  true
}

variable "vns3_license_type" {
  type = string
  default = "ul"
}

variable "vns3_instance_type" {
  default = "t3.medium"
}

variable "common_tags" {
  description = "A map of tags to add to all resources"
  default     = {
    ManagedBy = "Terraform"
    CreatedBy = "Cohesive solutions team"
    Topology  = "dev-env"
  }
}

variable "subnet_ids" {
    type = list
    default = [
        "subnet-06b9fc6a85df6e3c7"
    ]
}