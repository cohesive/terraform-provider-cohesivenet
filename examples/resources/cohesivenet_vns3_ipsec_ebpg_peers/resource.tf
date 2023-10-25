resource "cohesivenet_vns3_ipsec_ebpg_peers" "peer" {
  endpoint_id = vns3_ipsec_endpoints.endpoint.id
  ebgp_peer {
    ipaddress = "169.254.164.204"
    asn = 64512
    local_asn_alias = 65000
    access_list = "in permit 1.2.3.4/32, in permit 11.22.33.42/32, out permit 11.12.13.14/32"
    bgp_password = "password"
    add_network_distance = true
    add_network_distance_direction = "in"
    add_network_distance_hops = 10
  }
      depends_on = [
       vns3_ipsec_endpoints.endpoint
    ]
}