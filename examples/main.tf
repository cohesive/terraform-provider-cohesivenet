terraform {
  required_providers {
    cohesivenet = {
      version = "0.1.0"
      source  = "cohesive.net/vns3/cohesivenet"
    }
  }
}

provider "cohesivenet" {
  username = "vnscubed"
  password = "vnscontroller!"
  token = "771c844ecf0a2e0a9dd2c2a3071cfa7c1a06d7eed1f8664ce0995ec1b0824bee"
  hosturl = "https://3.127.171.216:8000/api"
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

data "cohesivenet_routes" routes {}

output "all_routes" {
   value = data.cohesivenet_routes.routes 
}