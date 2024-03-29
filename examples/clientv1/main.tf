terraform {
  required_providers {
    cohesivenet = {
      version = "0.1.0"
      source  = "cohesive.net/vns3/cohesivenet"
    }
  }
}

/* Configure Cohesive Terraform Provider

User this Terraform file is VNS3 controller already up and licensed

Specify host and either api_token or username/password */

provider "cohesivenet" {
  vns3 {
    username = "api"
    password = "mypassword"
    // api_token = "......"
    host = "116.170.38.40"
  }
}

/* Define datasources */

data "cohesivenet_vns3_route" all {}
data "cohesivenet_vns3_ipsec_endpoints" all {}
data "cohesivenet_vns3_config" config {}
data "cohesivenet_vns3_container_network" network {}
data "cohesivenet_vns3_firewall" all {}

/*  Define IPSec endpoint */
/*
resource "cohesivenet_vns3_ipsec_endpoints" "endpoint_1" {

  endpoint {
    name = "cohesive_to_peer"
    description = "cohesive_to_peer"
    ipaddress = "82.235.15.12"
    secret =  "verlongstring"
    pfs = true
    ike_version = 2
    nat_t_enabled = true
    extra_config = "phase1=aes256-sha2_256-dh16, phase2=aes256-sha2_256"
    vpn_type = "vti"
    route_based_int_address = "169.254.164.178/30"
    route_based_local =  "0.0.0.0/0"
    route_based_remote = "0.0.0.0/0"
  }
    
}
*/

/*  Add a VNS3 controller peer */

/*
resource "cohesivenet_vns3_ipsec_ebpg_peers" "peer" {

  endpoint_id = cohesivenet_vns3_ipsec_endpoints.endpoint_1.id
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
    cohesivenet_vns3_ipsec_endpoints.endpoint_1
  ]
}
*/

/*  Add a VNS3 controller route */

/* vns3 can be defined on a resource to override
   the default vns3 in provider
*/

/*
resource "cohesivenet_vns3_routes" "route" {
  vns3 {
    username = "api"
    password = "mypassword"
    host = "116.170.38.40"
  }
  route {
    cidr = "192.168.54.0/24"
    description = "cohesive_to_peer"
    interface = ""
    gateway = "192.168.54.0/24"
    advertise = true
    metric = 300
  }

}
*/

/* Example of how to use a map of routes */
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

/*  Add a VNS3 controller firewall rule */

/*
resource "cohesivenet_vns3_firewall_rules" "rule" {

  rule {
    script = "PREROUTING -d 10.18.0.65 -p udp --dport 162 -j DNAT --to 198.52.100.5:162"
  } 

}
*/


/*  Example of how to use a map of rules */

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

/*  Add a VNS3 controller plugin image */
/*
resource "cohesivenet_vns3_plugin_images" "image" {
  image {
    image_name = "test-tf-st-plugin"
    url  = "https://cohesive-networks-plugins.s3.amazonaws.com/plugins/vns3-high-availability/vns3-high-availability.pm.v2.2.1.tgz"
    //uildurl =
    //localbuild =
    //localimage =
    //imagefile =
    //buildfile =
    description = "test-tf-ha-description"
  }
 }

/*  Add a VNS3 controller plugin instance/

 resource  "cohesivenet_vns3_plugin_instances" instance {
    name = "pluginname"
    plugin_id = cohesivenet_vns3_plugin_images.image.id
    ip_address =  "198.51.100.11"
    description = "plugindescription"
    command = "/usr/bin/supervisord"
    environment = ""
    
    depends_on = [ cohesivenet_vns3_plugin_images.image ]
 }
*/

/*  Add a VNS3 controller web certificate */

/*
resource "cohesivenet_vns3_https_certs" "certs" {
  //filepath
  cert_file = var.vns3_license_cert_file
  key_file  = var.vns3_license_key_file
  //file
  cert = file("${path.module}/vns_cert.pem")
  key = file("${path.module}/vns_cert.key")
}
*/
