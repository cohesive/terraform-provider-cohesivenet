# Base Setup
This directory is a minimum install.  It deploys, licenses and configures two Peered VNS3 Controllers. 

## Starting the Deployment

Edit the `terraform.tfvars` file with the AWS ami id that is to be deployed.

Run `terraform init` to initialise the ternimal:
```bash
terraform init 
```

Run `terraform apply` to start the deployment, it will output the plan to the terminal. Type `yes` to run the deployment.
 ```bash
terraform apply 
```
Login information will be output after the deployment:  
vns3_instance_id_c1 = “i-xxxxxxxxxxxx”  
vns3_instance_id_c2 = “i-xxxxxxxxxxxx”  
vns3_public_ip_c1 = “x.x.x.x”  
vns3_public_ip_c2 = “x.x.x.x”  

Run `terraform destroy` to tear down the infrastructure.
```bash
terraform destroy 
``` 

## Expected Output

- 2 Licensed and Peered VNS3 controllers








