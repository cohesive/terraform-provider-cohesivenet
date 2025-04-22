  resource "cohesivenet_vns3_ipsec_endpoints" "endpoint" {
  endpoint {
      name = "endpoint_name"
      description = "endpoint_description"
      ipaddress = "123.123.123.123"
      secret =  "psk"
      pfs = true 
      ike_version = 2
      nat_t_enabled = true 
      private_ipaddress = "192.168.54.186"
      extra_config = "phase1=aes256-sha2_256-dh16, phase2=aes256-sha2_256-dh16"
      vpn_type = "vti"
      route_based_int_address = "169.254.32.186/30"
      route_based_local =  "0.0.0.0/0"
      route_based_remote = "0.0.0.0/0"
      }
        depends_on = [
          cohesivenet_vns3_ipsec_endpoints.endpoint
    ]
 }