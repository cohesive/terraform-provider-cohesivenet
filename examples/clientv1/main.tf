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
  token    = "771c844ecf0a2e0a9dd2c2a3071cfa7c1a06d7eed1f8664ce0995ec1b0824bee"
  host     = "https://3.127.171.216:8000/api"
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

data "data_source_vns3_container_network" all {}

output "all_container_networks" {
  value = data.cohesivenet_vns3_container_network.all 
  }


data "cohesivenet_vns3_routes" route {}

output "all_routes" {
   value = data.cohesivenet_vns3_routes.route 

}
data "cohesivenet_vns3_firewall" rules {}

output "all_rules" {
   value = data.cohesivenet_firewall.rules
}

*/
 resource "cohesivenet_vns3_ipsec_endpoints" "endpoint_vf" {
  endpoint {
      name = "cohesive_to_watford_secondary"
      description = "cohesive_to_watford_secondary245"
      ipaddress = "52.49.5.120"
      secret =  "OadQNYkfGB2R5UXpmv1mczlqgOTbpI8q"
      pfs = true
      ike_version = 2
      nat_t_enabled = true
      extra_config = "phase1=aes256-sha2_256-dh16 phase2=aes256-sha2_256"
      vpn_type = "vti"
      route_based_int_address = "169.254.164.178/30"
      route_based_local =  "10.18.0.64/26"
      route_based_remote = "0.0.0.0/0"
    }
    
 }
/*
 output "endpoint_vf_id" {
    value = cohesivenet_vns3_ipsec_endpoints.endpoint_vf.id
}

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
        depends_on = [
          cohesivenet_vns3_ipsec_endpoints.endpoint_vf
    ]
 }
 output "endpoint_vf2_id" {
    value = cohesivenet_vns3_ipsec_endpoints.endpoint_vf2.id
}

 resource "cohesivenet_vns3_ipsec_endpoints" "endpoint_vf3" {
  endpoint {
      name = "cohesive_to_watford_primary"
      description = "cohesive_to_watford_primary"
      ipaddress = "62.49.5.120"
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
         depends_on = [
          cohesivenet_vns3_ipsec_endpoints.endpoint_vf,
          cohesivenet_vns3_ipsec_endpoints.endpoint_vf2

    ]
    
 }

 output "endpoint_vf3_id" {
    value = cohesivenet_vns3_ipsec_endpoints.endpoint_vf3.id
}

  resource "cohesivenet_vns3_ipsec_endpoints" "endpoint_vf4" {
  endpoint {
      name = "cohesive_to_workload_primary"
      description = "cohesive_to_workload_primary"
      ipaddress = "44.241.205.244"
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
        depends_on = [
          cohesivenet_vns3_ipsec_endpoints.endpoint_vf2,
          cohesivenet_vns3_ipsec_endpoints.endpoint_vf3

    ]
 }
 output "endpoint_vf4_id" {
    value = cohesivenet_vns3_ipsec_endpoints.endpoint_vf4.id
}
*/
/*
resource "cohesivenet_vns3_ipsec_ebpg_peers" "peer" {
  endpoint_id = cohesivenet_vns3_ipsec_endpoints.endpoint_vf.id
  ebgp_peer {
    ipaddress = "169.254.164.204"
    asn = 64512
    local_asn_alias = 65000
    //access_list = "in permit any"
    //access_list = "in permit any, out permit any"
    access_list = "in permit 1.2.3.4/32, in permit 11.22.33.42/32, out permit 11.12.13.14/32"
    bgp_password = "password"
    add_network_distance = true
    add_network_distance_direction = "in"
    add_network_distance_hops = 10
  }
      depends_on = [
        cohesivenet_vns3_ipsec_endpoints.endpoint_vf
    ]
}
*/
/*

resource "cohesivenet_vns3_ipsec_ebpg_peers" "peer2" {
  endpoint_id = cohesivenet_vns3_ipsec_endpoints.endpoint_vf2.id
  ebgp_peer {
    ipaddress = "169.254.164.203"
    asn = 64512
    local_asn_alias = 65007
    access_list = "in permit any"
    bgp_password = "password"
    add_network_distance = true
    add_network_distance_direction = "in"
    add_network_distance_hops = 10
  }
    depends_on = [
        cohesivenet_vns3_ipsec_endpoints.endpoint_vf,
        cohesivenet_vns3_ipsec_endpoints.endpoint_vf2,
        cohesivenet_vns3_ipsec_endpoints.endpoint_vf3,
        cohesivenet_vns3_ipsec_endpoints.endpoint_vf4,
    ]
}

resource "cohesivenet_vns3_ipsec_ebpg_peers" "peer3" {
  endpoint_id = cohesivenet_vns3_ipsec_endpoints.endpoint_vf3.id
  ebgp_peer {
    ipaddress = "169.254.164.202"
    asn = 64512
    local_asn_alias = 65000
    access_list = "in permit 1.2.3.4/32,in permit 11.22.33.42/32,out permit 11.12.13.14/32"
    bgp_password = "password"
    add_network_distance = true
    add_network_distance_direction = "in"
    add_network_distance_hops = 10
  }
      depends_on = [
        cohesivenet_vns3_ipsec_endpoints.endpoint_vf,
        cohesivenet_vns3_ipsec_endpoints.endpoint_vf2,
        cohesivenet_vns3_ipsec_endpoints.endpoint_vf3,
        cohesivenet_vns3_ipsec_endpoints.endpoint_vf4,
    ]
}


resource "cohesivenet_vns3_ipsec_ebpg_peers" "peer4" {
  endpoint_id = cohesivenet_vns3_ipsec_endpoints.endpoint_vf4.id
  ebgp_peer {
    ipaddress = "169.254.164.201"
    asn = 64512
    local_asn_alias = 65007
    access_list = "in permit any"
    bgp_password = "password"
    add_network_distance = true
    add_network_distance_direction = "in"
    add_network_distance_hops = 10
  }
    depends_on = [
        cohesivenet_vns3_ipsec_endpoints.endpoint_vf,
        cohesivenet_vns3_ipsec_endpoints.endpoint_vf2,
        cohesivenet_vns3_ipsec_endpoints.endpoint_vf3,
        cohesivenet_vns3_ipsec_endpoints.endpoint_vf4,
    ]
}
*/



/*
 resource "cohesivenet_vns3_routes" "route" {
  route {
    cidr = "192.168.54.34/24"
    description = "cohesive_to_watford_secondary"
    interface = "tun0"
    gateway = "192.168.54.1/32"
    advertise = true
    metric = 300
  }
 }
*/
/*
variable "routes_map" {
  description = "Map of routes"
  type        = map(any)

  default = {
  "route": {
    "cidr": "192.168.54.102/32",
    "description": "cohesive_to_watford_one",
    "interface": "",
    "gateway": "",
    "advertise": true,
    "metric": 100
  },
  "route2": {
    "cidr": "192.168.54.112/32",
    "description": "cohesive_to_watford_two",
    "interface": "",
    "gateway": "",
    "advertise": true,
    "metric": 100
  },
  "route3": {
    "cidr": "192.168.54.113/32",
    "description": "cohesive_to_watford_three",
    "interface": "",
    "gateway": "",
    "advertise": true,
    "metric": 100
  }
}
}
*/
/*
variable "routes_map" {
  description = "Map of routes"
  type        = map(any)

  default = {
  "route": {
    "cidr": "192.168.54.10/32",
    "description": "cohesive_to_watford_secondary",
    "interface": "tun0",
    "gateway": "",
    "advertise": true,
    "metric": 100
  }
}
}
*/
/*
resource "cohesivenet_vns3_routes" "route-map" {
  dynamic route {
    for_each = var.routes_map
    content {
      cidr        = lookup(route.value, "cidr", null)
      description = lookup(route.value, "description", null)
      gateway     = lookup(route.value, "gateway", null)
      advertise   = lookup(route.value, "advertise", false)
    }
  }
}
*/
/*
routes_map = {
  "route": {
    "cidr": "192.168.54.10/32",
    "description": "cohesive_to_watford_secondary",
    "interface": "tun0",
    "gateway": "",
    "advertise": true,
    "metric": 100
  },
  "route2": {
    "cidr": "192.168.54.10/32",
    "description": "cohesive_to_watford_secondary",
    "interface": "tun0",
    "gateway": "",
    "advertise": true,
    "metric": 100
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
variable "rules_map" {
  description = "Map of rules"
  type        = map(any)

  default = {
  rule: {
    id = "10"
    script = "PREROUTING_CUST -d 10.18.0.65 -p udp --dport 162 -j DNAT --to 198.52.100.5:162"
  },
  rule1: {
    id = "21"
    script = "PREROUTING_CUST -d 10.18.0.66 -p udp --dport 162 -j DNAT --to 198.52.100.6:162"
  },
  rule2: {
    id = "31"
    script = "PREROUTING_CUST -d 10.18.0.67 -p udp --dport 162 -j DNAT --to 198.52.100.7:162"
  }
}
}

resource "cohesivenet_vns3_firewall_rules" "rule-map" {
  dynamic rule {
    for_each = var.rules_map
    content {
      id        = lookup(rule.value, "id", null)
      script = lookup(rule.value, "script", null)
    }
  }
}
*/
/*
variable "rules_map" {
  description = "Map of rules"
  type        = map(any)

  default = {
  rule: {
    script = "PREROUTING_CUST -d 10.18.0.65 -p udp --dport 162 -j DNAT --to 198.52.100.5:162"
  },
  rule1: {
    script = "PREROUTING_CUST -d 10.18.0.66 -p udp --dport 162 -j DNAT --to 198.52.100.6:162"
  },
  rule2: {
    script = "PREROUTING_CUST -d 10.18.0.67 -p udp --dport 162 -j DNAT --to 198.52.100.7:162"
  }
}
}

resource "cohesivenet_vns3_firewall_rules" "rule-map" {
  dynamic rule {
    for_each = var.rules_map
    content {
      script = lookup(rule.value, "script", null)
    }
  }
}


resource "cohesivenet_vns3_plugin_images" "image" {
  image {
    image_name = "test-tf-st-plugin"
    url  = "https://vns3-containers-read-all.s3.amazonaws.com/HA_Container/haplugin-pm.tar.gz"
    //uildurl =
    //localbuild =
    //localimage =
    //imagefile =
    //buildfile =
    description = "test-tf-ha-description"
  }
 }
 
 resource "cohesivenet_vns3_plugin_images" "nlb" {
  image {
    image_name = "test-tf-st-lb-plugin"
    url  = "https://st-temp-vf-share.s3.eu-central-1.amazonaws.com/nginx_lb_release_1509.export.tar.gz"
    //uildurl =
    //localbuild =
    //localimage =
    //imagefile =
    //buildfile =
    description = "test-tf-lb-description"
  }
 }

 resource  "cohesivenet_vns3_plugin_instances" instance {
    name = "pluginname"
    plugin_id = cohesivenet_vns3_plugin_images.image.id
    ip_address =  "198.51.100.11"
    description = "plugindescription"
    command = "/usr/bin/supervisord"
    environment = "HAENV_MODE=primary,HAENV_CLOUD=aws,HAENV_PEER_PUBLIC_IP=3.127.171.216,HAENV_SLEEP_TIME=15"
    
    //depends_on = [ cohesivenet_vns3_plugin_images.image ]
 }
 resource  "cohesivenet_vns3_plugin_instances" instance {
    name = "pluginname"
    plugin_id = cohesivenet_vns3_plugin_images.nlb.id
    ip_address =  "198.51.100.10"
    description = "plugindescription"
    command = "/usr/bin/supervisord"
    //environment = "server_list=[53,-,route-53-tcp],[53,udp,route-53-udp],[80,-,esm-workload-ingress]"
    environment = "server_list=\"[53,-,route-53-tcp],[53,udp,route-53-udp],[80,-,esm-workload-ingress]\",upstream_list=\"[route-53-tcp,10.23.2.43:53,10.23.2.59:53],[route-53-udp,10.23.2.43:53,10.23.2.59:53],[esm-workload-ingress,10.199.90.197:80,10.199.80.72:80]\""
    
    //depends_on = [ cohesivenet_vns3_plugin_images.image ]
 }
*/

/*
variable "vns3_license_cert_file" {
  # ADD PATH TO YOUR CERT FILE
  default = "/Users/scott/vns_cert.pem"
}
variable "vns3_license_key_file" {
  # ADD PATH TO YOUR KEY FILE
  default = "/Users/scott/vns_cert.key"
}

resource "cohesivenet_vns3_https_certs" "certs" {
  cert_file = var.vns3_license_cert_file
  key_file  = var.vns3_license_key_file
}
*/