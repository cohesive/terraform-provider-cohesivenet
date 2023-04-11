
output "vns3_public_ips_1" {
    value = aws_eip.vns3_ip_1.*.public_ip
}

output "vns3_public_ips_2" {
    value = aws_eip.vns3_ip_2.*.public_ip
}

output "vns3_private_dns_1" {
    value = aws_eip.vns3_ip_1.*.private_dns
}

output "vns3_private_dns_2" {
    value = aws_eip.vns3_ip_2.*.private_dns
}

output "vns3_1_topology_checksum" {
    value = cohesivenet_vns3_config.vns3_1.topology_checksum
}

output "vns3_1_keyset_checksum" {
    value = cohesivenet_vns3_config.vns3_1.keyset_checksum
}

output "vns3_1_peer_id" {
    value = cohesivenet_vns3_config.vns3_1.peer_id
}

output "vns3_2_topology_checksum" {
    value = cohesivenet_vns3_config.vns3_2.topology_checksum
}

output "vns3_2_keyset_checksum" {
    value = cohesivenet_vns3_config.vns3_2.keyset_checksum
}

output "vns3_2_peer_id" {
    value = cohesivenet_vns3_config.vns3_2.peer_id
}
