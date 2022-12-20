# Terraform Provider for Cohesive Networks VNS3 Controller

Beta version of Cohesive Networks Terraform provider for VNS3 cloud edge controller.

## Requirements
- Terraform .12.x+
- For building the provider, Install Go 1.18+


## Using Provider in Terraform Registry ( Terraform v0.12+ )

```
terraform {
  required_providers {
    cohesivenet = {
      version = "0.1.9"
      source  = "cohesive/cohesivenet"
    }
  }
}
```
### Quick Start

With a configured VNS3 Controller available, the following is a quick start example to add a route:

```
terraform {
  required_providers {
    cohesivenet = {
      source = "cohesive/cohesivenet"
      version = "0.1.9"
    }
  }
}

provider "cohesivenet" {
  vns3 {
    username = "vnscubed"
    password = "password"
    host = "host name or ip"
  }
}

resource "cohesivenet_vns3_routes" "route" {
  route {
    cidr = "192.168.54.0/24"
    description = "default route"
    advertise = true
    metric = 300
  }
 }
 ```

 1. Add the above snippet to a file called main.tf
 2. terraform init
 3. terraform apply

 See examples directory for setting up a single controller (basic) and peered controllers (config).

## Building the Provider

### Dependencies
This requires having go installed 1.18 or later installed on your machine and `go` in path. For Mac use brew with `brew install go`. You will also need [terraform installed](https://learn.hashicorp.com/tutorials/terraform/install-cli), v1.1 or later.

### Compiling provider
1. git clone this repository and cd into it
2. run `go mod vendor`
3. run `make install` - this will compile the provider (written in go) and move the binary to your local ~/.terraform.d/plugins directory so you can use the provider locally (see an example `terraform { required_providers {...` block in a examples/*/main.tf file). **Note** - you may need to change the architecture used in the [makefile](./Makefile). It defaults to darwin_amd64. M1 will require the ARM arch. you will change it under the `install` command.

Explanation:
- `source = "cohesive.net/vns3/cohesivenet"` - this is defined in the Makefile with **HOSTNAME**, **NAMESPACE** and **NAME** variables.
- this is an example that also uses an aws provider

**Running Terraform**

1. cd into your terraform dir. There are examples under `./examples`. e.g. cd `examples/config`.  Note see this snippet for using local built provider
```
cohesivenet = {
      version = "0.1.0"
      source  = "cohesive.net/vns3/cohesivenet"
    }
```
2. `terraform init` - this will initialize your state and install plugins, including the locally installed cohesivenet provider
3. `terraform plan` - review the resources that will be created. *hint* - its good practice to save these plans and them build them directly with `-out`. e.g. `terraform plan -out "my-plan-$(date -u +"%Y-%m-%dT%H-%M-%SZ").tfplan"` will output a timestamped plan
4. Build with `terraform apply`. if you have a plan from step 3, you can pass that as first argument: `terraform apply my-plan-2022-01-01T10:31:03Z.tfplan`

**Upgrading plugin for your terraform code**

Remove the lock file and run a terraform upgrade
`rm .terraform.lock.hcl && terraform init -upgrade`

When you re-run `make install` this will update the plugin installed locally but it will not automatically update the plugin used by your terraform code. For instance, when testing, you run `make install` and then cd into examples/config and here you will need to pull in your updated plugin code. When you ran #2 above `terraform init`, terraform downloaded the plugins defined in main.tf and then wrote a file called [.terraform.lock.hcl](https://learn.hashicorp.com/tutorials/terraform/provider-versioning#explore-terraform-lock-hcl) to the same directly. This lock file ensures other users of the same infra code use the correct provider revision. To upgrade we need you will need to remove the lock file and run `terraform init -upgrade`.


