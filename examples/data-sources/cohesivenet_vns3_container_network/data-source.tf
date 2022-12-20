data "cohesivenet_vns3_container_network" network {}

output "all_container_networks" {
  value = data.cohesivenet_vns3_container_network.network 
}