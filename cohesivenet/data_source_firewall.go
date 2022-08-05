package cohesivenet

import (
	"context"
	"strconv"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFirewall() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFirewallRead,
		Schema: map[string]*schema.Schema{
			"rules": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceFirewallRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	firewallResponse, err := c.GetFirewallRules()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("rules", flattenRules(firewallResponse)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenRules(firewallResponse cn.FirewallResponse) interface{} {
	rules := make([]interface{}, len(firewallResponse.FirewallRules), len(firewallResponse.FirewallRules))
	i := 0
	for _, rt := range firewallResponse.FirewallRules {
		row := make(map[string]interface{})
		row["id"] = rt.ID
		row["rule"] = rt.Rule
		rules[i] = row
		i++
	}
	return rules
}
