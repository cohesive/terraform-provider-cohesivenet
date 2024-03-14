
resource "cohesivenet_vns3_license_upgrade" "upgrade" {
  license_upgrade_key = "/path/to/license/file"
  clientpack_ips = "100.127.255.26,100.127.255.27,100.127.255.28,100.127.255.29,100.127.255.30"
  manager_ips = "100.127.255.252,100.127.255.253"
}