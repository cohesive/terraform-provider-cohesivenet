data "cohesivenet_vns3_firewall" all {}

output "all_rules" {
   value = data.cohesivenet_vns3_firewall.all
}