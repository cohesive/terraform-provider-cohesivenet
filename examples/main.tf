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
}

//data "cohesivenet_endpoints" all {}

///output "all_endpoints" {
//   value = data.cohesivenet_endpoints.all
//}


data "cohesivenet_config" all {}

output "all_config" {
  value = data.cohesivenet_config.all 
}

data "cohesivenet_container_network" all {}

output "all_container_networks" {
  value = data.cohesivenet_container_network.all 
}

//data "cohesivenet_routes" routes {}

//output "all_routes" {
//   value = data.cohesivenet_routes.routes 
//}