terraform {
  required_providers {
    cohesivenet = {
      version = "0.1.0"
      source  = "cohesive.net/vns3/cohesivenet"
    }
  }
}

provider "cohesivenet" {}

resource "cohesivenet_vns3_link" "linkA" {
    vns3 {
        host = var.vns3_host
        password = var.vns3_password
        timeout = 60 // creating a link can sometimes take longer
    }

    link_id = var.link_id
    name = var.link_name
    description = var.link_description
    conf = file(var.link_conf_file)
    policies = try(file(var.link_policies_file), null)
}

output "link_ip" {
    value = cohesivenet_vns3_link.linkA.clientpack_ip
}

output "link_id" {
    value = cohesivenet_vns3_link.linkA.id
}

output "link_type" {
    value = cohesivenet_vns3_link.linkA.type
}

output "link_name" {
    value = cohesivenet_vns3_link.linkA.name
}
