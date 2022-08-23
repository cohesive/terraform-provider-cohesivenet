# Configuration example

This uses the following resources from the cohesivenet provider:

1. `cohesivenet_vns3_config` - this resource will license a controller, generate a keyset, setup peering and optionall generate new authentication credentials. It can also be used to configure a VNS3 controller by fetching the configuration from a peer VNS3
2. `cohesivenet_vns3_peers` - this configures a VNS3 controllers peers

**Note: There is a natural dependency between the AWS instance resources and the vns3 config resources as the controller Id and aws_eip outputs are passed to the vns3 config resources. But the cohesivenet_vns3_peers should use a depends_on for explicitness**


## Variables to change

**vns3_master_password** - used for UI and API

**vns3_license_file** - path to a license file on our computer

**subnet_ids** - subnets to launch VNS3s inside

**security_group_id** - security group to launch VNS3 inside
