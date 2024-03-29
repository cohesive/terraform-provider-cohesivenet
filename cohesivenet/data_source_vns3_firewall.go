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
			"rule": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of rule in firewall",
						},
						"script": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Firewall rule in VNS3 syntax",
						},
					},
				},
			},
		},
	}
}

func dataSourceFirewallRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	var diags diag.Diagnostics

	firewallResponse, err := c.GetFirewallRules()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("rule", flattenRules(firewallResponse)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenRules(firewallResponse cn.FirewallResponse) interface{} {
	rules := make([]interface{}, len(firewallResponse.FirewallRules))
	i := 0
	for _, rt := range firewallResponse.FirewallRules {
		row := make(map[string]interface{})
		row["id"] = rt.ID
		row["script"] = rt.Rule
		rules[i] = row
		i++
	}
	return rules
}
