resource "vns3_firewall_rules" "rule" {
    script = "PREROUTING_CUST -d 10.10.10.10 -p udp --dport 123 -j DNAT --to 192.168.1.1:123"
    }