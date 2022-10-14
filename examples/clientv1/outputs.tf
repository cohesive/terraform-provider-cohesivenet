output "all_endpoints" {
   value = data.cohesivenet_vns3_ipsec_endpoints.all
}

output "all_config" {
  value = data.cohesivenet_vns3_config.config
}

output "all_container_networks" {
  value = data.cohesivenet_vns3_container_network.all 
  }

output "all_routes" {
   value = data.cohesivenet_vns3_routes.route 

}

output "all_rules" {
   value = data.cohesivenet_firewall.rules
}

 output "endpoint_1_id" {
    value = cohesivenet_vns3_ipsec_endpoints.endpoint_1_id.id
}