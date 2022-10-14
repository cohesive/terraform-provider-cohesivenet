 resource "vns3_routes" "route" {
  route {
    cidr = "192.168.54.34/24"
    description = "route_description"
    interface = "tun0"
    gateway = "192.168.54.1/32"
    advertise = true
    metric = 300
  }
 }