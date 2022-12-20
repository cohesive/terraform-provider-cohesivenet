data "cohesivenet_vns3_config" config {}

output "all_config" {
  value = data.cohesivenet_vns3_config.config
}