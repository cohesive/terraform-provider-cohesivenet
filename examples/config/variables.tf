
variable "vns3_master_password" {
  sensitive = true
  type = string
}

variable "vns3_api_token_lifetime" {
  type = number
  description = "time in seconds until generated VNS3 token expires. 0 means don't generate token"
  default = 0
}

variable "vns3_api_token_refresh" {
  type = bool
  description = "token lifetime will refresh for each successful request"
  default = false
}

variable "topology_name" {
    default = "bens-test-tf-launch"
}
variable "controller_name" {
    default = "ctrl"
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
  default = "t3.small"
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
        "subnet-06b9fc6a85df6e3c7",
        "subnet-02cea8256e7515eca"
    ]
}

variable "security_group_id" {
    default = "sg-011082922b0f4915f"
}