terraform {
  required_providers {
    cohesivenet = {
      version = "0.1.0"
      source  = "cohesive.net/vns3/cohesivenet"
    }
  }
}

provider "cohesivenet" {
  username = ""
  password = ""
  token = ""
  hosturl = ""
}

data "cohesivenet_endpoints" all {}

output "all_endpoints" {
   value = data.cohesivenet_endpoints.all
}

/*
data "cohesivenet_config" config {}

output "all_config" {
  value = data.cohesivenet_config.config
}


resource "routes" "1" {
    		cidr = "x.x.x.x"
			  description = "description"
}

//data "cohesivenet_container_network" all {}

//output "all_container_networks" {
//  value = data.cohesivenet_container_network.all 
//}

data "cohesivenet_routes" route {}

output "all_routes" {
   value = data.cohesivenet_routes.route 
}
data "cohesivenet_firewall" rules {}

output "all_rules" {
   value = data.cohesivenet_firewall.rules
}
*/

/*
 resource "cohesivenet_endpoints" "endpoint_vf" {
  endpoint {
      endpoint_name = "route_based_vf"
      description = "routebased_api"
      peer_ip = "3.64.150.23"
      secret =  "biglongstring"
      pfs = true
      ike_version= 2
      nat_t_enabled = true
      extra_config = "phase1=aes256-sha1-dh14"
      vpn_type = "vti"
      route_based_int_address = "169.254.0.70/30"
      route_based_local =  "0.0.0.0/0"
      route_based_remote = "0.0.0.0/0"
    }
 }
*/



