resource "cohesivenet_vns3_firewall_rule" "rule" {
    position = 1
    rule = "PREROUTING_CUST -d 10.10.10.10 -p udp --dport 123 -j DNAT --to 192.168.1.1:123"
}