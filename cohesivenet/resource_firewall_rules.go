package cohesivenet

import (
	"context"

	cn "github.com/cohesive/cohesivenet-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRulesCreate,
		ReadContext:   resourceRulesRead,
		UpdateContext: resourceRulesUpdate,
		DeleteContext: resourceRulesDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rules": &schema.Schema{
				Type:     schema.TypeList,
				ForceNew: true,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"rule": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceRulesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cn.Client)

	var diags diag.Diagnostics

	rle := d.Get("rules").([]interface{})[0]
	rule := rle.(map[string]interface{})

	rl := cn.FirewallRule{
		Rule: rule["rule"].(string),
	}

	newRule, err := c.CreateFirewallRules(&rl)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newRule.FirewallRules[0].ID)

	//resourceRoutesRead(ctx, d, m)

	return diags
}

func resourceRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}

func resourceRulesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceEndpointsRead(ctx, d, m)
}

func resourceRulesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cn.Client)

	var diags diag.Diagnostics

	ruleId := d.Id()

	err := c.DeleteRule(ruleId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func flattenRulesData(ruleResponse cn.FirewallResponse) []interface{} {
	routes := make([]interface{}, len(ruleResponse.FirewallRules), len(ruleResponse.FirewallRules))

	i := 0
	for _, rt := range ruleResponse.FirewallRules {
		row := make(map[string]interface{})

		row["id"] = rt.ID
		row["rule"] = rt.Rule

		routes[i] = row
		i++
	}

	return routes

}
