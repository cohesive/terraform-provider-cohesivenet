
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


variable "vns3_license_cert_file" {
  # ADD PATH TO YOUR CERT FILE
  default = "/path/to/vns_cert.pem"
}
variable "vns3_license_key_file" {
  # ADD PATH TO YOUR KEY FILE
  default = "/path/to/vns_cert.key"
}