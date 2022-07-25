package cohesivenet

import (
	"context"
	"strconv"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFirewall() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFirewallRead,
		Schema: map[string]*schema.Schema{
			"response": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
	c := m.(*cn.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	firewallResponse, err := c.GetFirewallRules()
	if err != nil {
		return diag.FromErr(err)
	}

	//rules := firewallResponse.Firewall

	/*
		rulesMap := make(map[string]string)

		i := 0
		for _, rl := range rules {
			row := make(map[string]interface{})
			//rt_data := rl.(interface{})
			row["rule"] = rl

			rulesMap[] = row
			i++
		}
	*/
	//	rules := make([]interface{}, 1, 1)
	rules := firewallResponse.Firewall
	//row := make(map[string]interface{})
	//row["rule"] = firewallResponse.Firewall[0]

	if err := d.Set("response", rules); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

/*
func flattenRules(routeResponse map[string]interface{}) interface{} {
	if routeResponse != nil {
		routes := make([]interface{}, len(routeResponse), len(routeResponse))

		i := 0
		for id, rt := range routeResponse {
			row := make(map[string]interface{})
			rt_data := rt.(map[string]interface{})

			row["cidr"] = rt_data["cidr"]
			row["id"] = id
			row["description"] = rt_data["description"]
			row["advertise"] = rt_data["advertise"].(bool)
			row["metric"] = rt_data["metric"].(float64)
			row["enabled"] = rt_data["enabled"].(bool)
			row["netmask"] = rt_data["netmask"]
			row["editable"] = rt_data["editable"].(bool)
			row["table"] = rt_data["table"]
			row["interface"] = rt_data["interface"]

			routes[i] = row
			i++
		}

		return routes
	}

	return make([]interface{}, 0)
}
*/
