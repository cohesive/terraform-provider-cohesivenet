
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
    # CHANGE TOPOLOGY NAME
    default = "terraform-launch"
}
variable "controller_name" {
    default = "terraform"
}

variable "vns3_account_owner" {
  type = string
  #default = "201138274120"
  default = "678344834199"
}

variable "vns3_version" {
    type = string
    default = "602-20221104"
}

variable "vns3_license_file" {
    # ADD PATH TO YOUR LICENSE FILE
  default = "~/license.txt"
}

variable "keyset_token" {
    default = "keysetpassword"
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
  # CHANGE subnet ids
    type = list
    default = [
        "subnet-07823348e19526c990",
        "subnet-0de18940602c1d1343"
    ]
    description = "The subnets to launch VNS3 controllers. 1 VNS3 controller for each subnet"
}

variable "security_group_id" {
  # CHANGE security group
  default = "sg-055a3082672c469825"
  description = "The Security group to launch VNS3 controllers in"
}
