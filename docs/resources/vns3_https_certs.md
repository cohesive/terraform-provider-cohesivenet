---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cohesivenet_vns3_https_certs Resource - terraform-provider-cohesivenet"
subcategory: ""
description: |-
  Uploads HTTPS certificates to the VNS3 controller UI.
---

# cohesivenet_vns3_https_certs (Resource)

Uploads HTTPS certificates to the VNS3 controller UI.

## Example Usage
```terraform
resource "vns3_https_certs" "cert" {
  cert_file = "/Path/to/vns_cert.pem"
  key_file = "/Path/to/vns_cert.key"
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `cert_file` (String) File path to certificate.
- `key_file` (String) File path to private key.

### Read-Only

- `id` (String) The ID of this resource.

