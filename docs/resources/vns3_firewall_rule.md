---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cohesivenet_vns3_firewall_rule Resource - terraform-provider-cohesivenet"
subcategory: ""
description: |-
  Creates firewall rule using the Cohesive simplified IP tables syntax.
---

# cohesivenet_vns3_firewall_rule (Resource)

Creates firewall rule using the Cohesive simplified IP tables syntax.

## Example Usage

```terraform
resource "cohesivenet_vns3_firewall_rule" "rule" {
    position = 1
    rule = "PREROUTING_CUST -d 10.10.10.10 -p udp --dport 123 -j DNAT --to 192.168.1.1:123"
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `rule` (String) Firewall rule in VNS3 syntax

### Optional

- `comment` (String) Firewall comment
- `disabled` (Boolean) Whether the rule is disabled or not
- `groups` (List of String) Firewall rule groups
- `position` (Number) Position of firewall rule

### Read-Only

- `created_at` (String) Firewall rule created date
- `id` (String) The ID of this resource.
- `last_resolved` (String) Firewall rule date last resolved
- `rule_resolved` (String) Firewall rule resolved table in VNS3 syntax
- `table` (String) Firewall rule table in VNS3
- `last_updated` (String)


