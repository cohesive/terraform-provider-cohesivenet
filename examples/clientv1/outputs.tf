

output "all_routes" {
   value = data.cohesivenet_vns3_route.all 
}


output "all_config" {
  value = data.cohesivenet_vns3_config.config
}


output "all_container_networks" {
  value = data.cohesivenet_vns3_container_network.network 
  }


output "all_rules" {
   value = data.cohesivenet_vns3_firewall.all
}


 output "all_endpoints" {
    value = data.cohesivenet_vns3_ipsec_endpoints.all
}
