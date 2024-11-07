output "vns3_public_ip_c1" {
  value = aws_eip.vns3_ip_1.public_ip
}

output "vns3_instance_id_c1" {
  value = aws_instance.vns3controller_1.id
}
