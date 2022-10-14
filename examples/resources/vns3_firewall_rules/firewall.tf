resource "vns3_firewall_rules" "rule" {
    script = "PREROUTING_CUST -d 10.18.0.65 -p udp --dport 162 -j DNAT --to 198.52.100.5:162"
    }