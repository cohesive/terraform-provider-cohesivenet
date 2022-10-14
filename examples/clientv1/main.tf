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


data "cohesivenet_vns3_ipsec_endpoints" all {}

data "cohesivenet_vns3_config" config {}

data "data_source_vns3_container_network" all {}

data "cohesivenet_vns3_routes" route {}

data "cohesivenet_vns3_firewall" rules {}

output "all_rules" {
   value = data.cohesivenet_firewall.rules
}

 resource "cohesivenet_vns3_ipsec_endpoints" "endpoint_1" {
  endpoint {
      name = "cohesive_to_peer"
      description = "cohesive_to_peer"
      ipaddress = "82.235.15.12"
      secret =  "verlongstring"
      pfs = true
      ike_version = 2
      nat_t_enabled = true
      extra_config = "phase1=aes256-sha2_256-dh16 phase2=aes256-sha2_256"
      vpn_type = "vti"
      route_based_int_address = "169.254.164.178/30"
      route_based_local =  "0.0.0.0/0"
      route_based_remote = "0.0.0.0/0"
    }
    
 }

resource "cohesivenet_vns3_ipsec_ebpg_peers" "peer" {
  endpoint_id = cohesivenet_vns3_ipsec_endpoints.endpoint_vf.id
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
        cohesivenet_vns3_ipsec_endpoints.endpoint_vf
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


resource "cohesivenet_firewalls" "rule" {
  rule {
    script = "PREROUTING -d 10.18.0.65 -p udp --dport 162 -j DNAT --to 198.52.100.5:162"
  }
  
}


/*
resource "cohesivenet_vns3_firewall_rules" "rule-map" {
  dynamic rule {
    for_each = var.rules_map
    content {
      script = lookup(rule.value, "script", null)
    }
  }
}
*/
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


 resource  "cohesivenet_vns3_plugin_instances" instance {
    name = "pluginname"
    plugin_id = cohesivenet_vns3_plugin_images.image.id
    ip_address =  "198.51.100.11"
    description = "plugindescription"
    command = "/usr/bin/supervisord"
    environment = "HAENV_MODE=primary,HAENV_CLOUD=aws,HAENV_PEER_PUBLIC_IP=3.127.171.216,HAENV_SLEEP_TIME=15"
    
    //depends_on = [ cohesivenet_vns3_plugin_images.image ]
 }

resource "cohesivenet_vns3_https_certs" "certs" {
  cert_file = var.vns3_license_cert_file
  key_file  = var.vns3_license_key_file
}