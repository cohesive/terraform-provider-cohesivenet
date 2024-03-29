---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cohesivenet_vns3_ipsec_ebpg_peers Resource - terraform-provider-cohesivenet"
subcategory: ""
description: |-
  Creates eBGP peer in conjunction with vns3_ipsec_endpoints resource.
---

# cohesivenet_vns3_ipsec_ebpg_peers (Resource)

Creates eBGP peer in conjunction with vns3_ipsec_endpoints resource.

## Example Usage

```terraform
resource "cohesivenet_vns3_ipsec_ebpg_peers" "peer" {
  endpoint_id = vns3_ipsec_endpoints.endpoint.id
  ebgp_peer {
    ipaddress = "169.254.164.204"
    asn = 64512
    local_asn_alias = 65000
    access_list = "in permit 1.2.3.4/32, in permit 11.22.33.42/32, out permit 11.12.13.14/32"
    bgp_password = "password"
    add_network_distance = true
    add_network_distance_direction = "in"
    add_network_distance_hops = 10
  }
      depends_on = [
       vns3_ipsec_endpoints.endpoint
    ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `ebgp_peer` (Block List, Min: 1) Nested block for eBGP peer attributes (see [below for nested schema](#nestedblock--ebgp_peer))
- `endpoint_id` (Number)

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--ebgp_peer"></a>
### Nested Schema for `ebgp_peer`

Required:

- `asn` (Number) Autonomous System Number of your network
- `ipaddress` (String) IP address or neighbor IP for BGP

Optional:

- `access_list` (String) Access Control List. IN PERMIT xxxx / OUT PERMIT xxxx
- `add_network_distance` (Boolean) Specifies if we are using network distance weighting, Default: false
- `add_network_distance_direction` (String) Specifies direction for distance weighting. IN / OUT
- `add_network_distance_hops` (Number) Specifies how many hops for network distance weighting
- `bgp_password` (String) Password for BGP, if required
- `local_asn_alias` (Number) ASN alias

Read-Only:

- `id` (String) Id of the eBGP peer


