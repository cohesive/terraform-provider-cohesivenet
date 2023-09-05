resource "cohesivenet_vns3_ipsec_traffic_pair" "traffic_pair" {
  provider = cohesivenet.controller_1
  endpoint_id = cohesivenet_vns3_ipsec_endpoints.endpoint.id
  remote_subnet = "172.31.11.20/32"
  local_subnet = "10.10.0.42/32"
  description = "test"
  ping_interval = 30
	ping_interface = "eth0"
  ping_ipaddress = "10.10.0.36"
	enabled = true

  depends_on = [
        cohesivenet_vns3_ipsec_endpoints.endpoint
    ]
}
