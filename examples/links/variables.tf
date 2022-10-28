
variable "vns3_password" {
  sensitive = true
  type = string
}

variable "vns3_host" {
  type = string
}

variable "link_id" {
    type = number
}

variable "link_name" {
    type = string
}

variable "link_description" {
  type = string
  default = ""
}

variable "link_conf_file" {
  type = string
  description = "Link clientpack conf file"
}

variable "link_policies_file" {
  type = string
  description = "Extra polices to add to end of clientpack"
  default = ""
}

