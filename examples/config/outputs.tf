output "vns3_public_ip" {
    value = aws_eip.vns3_ip.public_ip
}

output "vns3_topology_checksum" {
    value = cohesivenet_vns3_config.vns3.topology_checksum
}

output "vns3_keyset_checksum" {
    value = cohesivenet_vns3_config.vns3.keyset_checksum
}

output "vns3_peer_id" {
    value = cohesivenet_vns3_config.vns3.peer_id
}