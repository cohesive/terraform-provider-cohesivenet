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
data "cohesivenet_vns3_ipsec_endpoints" all {}

output "all_endpoints" {
   value = data.cohesivenet_vns3_ipsec_endpoints.all
}

data "cohesivenet_vns3_config" config {}

output "all_config" {
  value = data.cohesivenet_vns3_config.config
}

data "cohesivenet_container_network" all {}

output "all_container_networks" {
  value = data.cohesivenet_vns3_container_network.all }

data "cohesivenet_vns3_routes" route {}

output "all_routes" {
   value = data.cohesivenet_vns3_routes.route 

}
data "cohesivenet_vns3_firewall" rules {}

output "all_rules" {
   value = data.cohesivenet_firewall.rules
}
*/
/*
 resource "cohesivenet_vns3_ipsec_endpoints" "endpoint_vf" {
  endpoint {
      name = "cohesive_to_watford_secondary"
      description = "cohesive_to_watford_secondary245"
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
*/
/*
  resource "cohesivenet_vns3_ipsec_endpoints" "endpoint_vf2" {
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
*/

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
resource "cohesivenet_firewalls" "rule" {
  rule {
    id = "0"
    script = "PREROUTING -d 10.18.0.65 -p udp --dport 162 -j DNAT --to 198.52.100.5:162"
  }
  
}

resource "cohesivenet_firewalls" "rule1" {
  rule {
    id = "1"
    script = "PREROUTING -d 10.18.0.66 -p udp --dport 162 -j DNAT --to 198.52.100.6:162"
  }
  
}
*/

/*
resource "cohesivenet_vns3_ipsec_ebpg_peers" "peer" {
  endpoint_id = 1
  ebgp_peer {
    ipaddress = "169.254.164.177"
    asn = 64512
    local_asn_alias = 65000
    access_list = "in permit any"
    bgp_password = "password"
    add_network_distance = true
    add_network_distance_direction = "in"
    add_network_distance_hops = 10
  }
      depends_on = [
        cohesivenet_ipsec_endpoints.endpoint_vf
    ]
}
*/
/*
resource "cohesivenet_vns3_ipsec_ebpg_peers" "peer2" {
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
    depends_on = [
        cohesivenet_ipsec_endpoints.endpoint_vf2
    ]
}
*/


resource "cohesivenet_vns3_plugin_images" "image" {
  image {
    name = "test-tf-st-plugin"
    url  = "https://vns3-containers-read-all.s3.amazonaws.com/HA_Container/haplugin-pm.tar.gz"
    //uildurl =
    //localbuild =
    //localimage =
    //imagefile =
    //buildfile =
    description = "test-tf-ha-description"
  }
 }

 resource  "cohesivenet_vns3_plugin_instances" instance {
    name = "pluginname"
    //plugin_id = "sha256:9fe7429af80c9d1a8d53aa4f16f72bde0c73e153783cbdf0c95b23917b428e83" // var of cohesivenet_vns3_plugin_images.image.id
    plugin_id = cohesivenet_vns3_plugin_images.image.image[0].id
    ip_address =  "198.51.100.11"
    description = "plugindescription"
    command = "/usr/bin/supervisord"
    
    depends_on = [
    cohesivenet_vns3_plugin_images.image
    ]
 }

