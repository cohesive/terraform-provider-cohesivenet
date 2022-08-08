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
  host = "https://3.127.171.216:8000/api"
}

/*
data "cohesivenet_endpoints" all {}

output "all_endpoints" {
   value = data.cohesivenet_endpoints.all
}

data "cohesivenet_config" config {}

output "all_config" {
  value = data.cohesivenet_config.config
}

data "cohesivenet_container_network" all {}

output "all_container_networks" {
  value = data.cohesivenet_container_network.all }

data "cohesivenet_routes" route {}

output "all_routes" {
   value = data.cohesivenet_routes.route 

}
data "cohesivenet_firewall" rules {}

output "all_rules" {
   value = data.cohesivenet_firewall.rules
}
*/

 resource "cohesivenet_endpoints" "endpoint_vf" {
  endpoint {
      name = "cohesive_to_watford_secondary"
      description = "cohesive_to_watford_secondary"
      ipaddress = "52.49.5.120"
      secret =  "OadQNYkfGB2R5UXpmv1mczlqgOTbpI8q"
      pfs = true
      ike_version = 2
      nat_t_enabled = true
      extra_config = "phase1=aes256-sha2_256-dh16"
      vpn_type = "vti"
      route_based_int_address = "169.254.164.178/30"
      route_based_local =  "10.18.0.64/26"
      route_based_remote = "0.0.0.0/0"
    }
 }


  resource "cohesivenet_endpoints" "endpoint_vf2" {
  endpoint {
      name = "cohesive_to_workload_secondary"
      description = "cohesive_to_workload_secondary"
      ipaddress = "34.241.205.244"
      secret =  "wakeYgPcjh6IJ70NmLxzrkSE0Wz9guMP"
      pfs = true
      ike_version = 2
      nat_t_enabled = true
      extra_config = "phase1=aes256-sha2_256-dh16"
      vpn_type = "vti"
      route_based_int_address = "169.254.32.186/30"
      route_based_local =  "0.0.0.0/0"
      route_based_remote = "0.0.0.0/0"
    }
 }


/*
 resource "cohesivenet_routes" "route" {
  route {
    cidr = "192.168.54.0/24"
    description = "cohesive_to_watford_secondary"
    interface = "tun0"
    gateway = "192.168.54.1/32"
    advertise = true
    metric = 300
  }
 }
*/
/*
resource "cohesivenet_firewall" "rule" {
  rules {
    id = "0"
    rule = "PREROUTING -d 10.18.0.65 -p udp --dport 162 -j DNAT --to 198.52.100.5:162"
  }
  
}

resource "cohesivenet_firewall" "rule1" {
  rules {
    id = "1"
    rule = "PREROUTING -d 10.18.0.66 -p udp --dport 162 -j DNAT --to 198.52.100.6:162"
  }
  
}
*/


resource "cohesivenet_ipsec_ebpg" "peer" {
  endpoint_id = 1
  ebgp_peer {
    ipaddress = "169.254.164.177"
    asn = 64512
    local_asn_alias = 65007
    access_list = "in permit any"
    bgp_password = "password"
    add_network_distance = true
    add_network_distance_direction = "in"
    add_network_distance_hops = 10
  }
}
resource "cohesivenet_ipsec_ebpg" "peer2" {
  endpoint_id = 2
  ebgp_peer {
    ipaddress = "169.254.164.178"
    asn = 64512
    local_asn_alias = 65007
    access_list = "in permit any"
    bgp_password = "password"
    add_network_distance = true
    add_network_distance_direction = "in"
    add_network_distance_hops = 10
  }
}
