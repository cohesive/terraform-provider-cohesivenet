data "cohesivenet_vns3_route" all {}

output "all_routes" {
   value = data.cohesivenet_vns3_route.all 
}