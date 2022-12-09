data "cohesivenet_vns3_ipsec_endpoints" all {}

output "all_endpoints" {
    value = data.cohesivenet_vns3_ipsec_endpoints.all
}
