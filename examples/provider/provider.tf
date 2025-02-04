terraform {
  required_providers {
    cohesivenet = {
      source = "cohesive/cohesivenet"
      version = "0.1.8"
    }
  }
}


# Configure Cohesive Terraform Provider Specify host and either api_token or username/password.
# For multiple controllers a Provider Alias can be used. 


provider "cohesivenet" {
  alias = "controller_1"
  username = "default_username"
  host = aws_eip.vns3_ip_1.public_ip
  password = aws_instance.vns3controller_1.id
  //token = cohesivenet_vns3_config.vns3_1.token
}


provider "cohesivenet" {
  alias = "controller_2"
  username = "default_username"
  host = aws_eip.vns3_ip_2.public_ip
  password = aws_instance.vns3controller_2.id
  //token = cohesivenet_vns3_config.vns3_2.token
}