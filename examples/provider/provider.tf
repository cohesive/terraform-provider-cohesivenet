terraform {
  required_providers {
    cohesivenet = {
      source = "cohesive/cohesivenet"
      version = "0.1.8"
    }
  }
}


# Configure Cohesive Terraform Provider Specify host and either api_token or username/password

provider "cohesivenet" {
  vns3 {
    username = "api"
    password = "mypassword"
    // api_token = "......"
    host = "11.22.33.44"
  }
}