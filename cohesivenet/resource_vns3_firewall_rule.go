package cohesivenet

import (
	"context"
	"strconv"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet/v1"
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
			"rule": &schema.Schema{
				Type:        schema.TypeList,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "Nested Block for rules",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							ForceNew:    true,
							Optional:    true,
							Computed:    true,
							Description: "Id given to rule after it has been applied",
						},
						"script": &schema.Schema{
							Type:        schema.TypeString,
							ForceNew:    true,
							Optional:    true,
							Description: "Firewall rule in VNS3 syntax",
						},
					},
				},
			},
		},
	}
}

func resourceRulesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	var diags diag.Diagnostics
	var ruleList []*cn.FirewallRule
	rules := d.Get("rule").([]interface{})

	for _, rule := range rules {
		rle := rule.(map[string]interface{})
		rule := cn.FirewallRule{
			Rule: rle["script"].(string),
		}

		ruleList = append(ruleList, &rule)
	}
	err := c.CreateFirewallRules(ruleList)
	if err != nil {
		return diag.FromErr(err)
	}
	//d.SetId(newRule.FirewallRules[0].ID)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	resourceRulesRead(ctx, d, m)

	return diags
}

/*
func resourceRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}
*/

func resourceRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	firewallResponse, err := c.GetFirewallRules()
	if err != nil {
		return diag.FromErr(err)
	}

	rules := flattenRulesData(firewallResponse)

	d.Set("rule", rules)

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func resourceRulesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceEndpointsRead(ctx, d, m)
}

func resourceRulesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	var diags diag.Diagnostics

	var ruleList []*cn.FirewallRule
	rules := d.Get("rule").([]interface{})

	for _, rule := range rules {
		rle := rule.(map[string]interface{})
		rule := cn.FirewallRule{
			ID:   rle["id"].(string),
			Rule: rle["script"].(string),
		}

		ruleList = append(ruleList, &rule)
	}

	err := c.DeleteRules(ruleList)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func flattenRulesData(ruleResponse cn.FirewallResponse) []interface{} {
	//routes := make([]interface{}, len(ruleResponse.FirewallRules), len(ruleResponse.FirewallRules))
	routes := make([]interface{}, len(ruleResponse.FirewallRules))

	i := 0
	for _, rt := range ruleResponse.FirewallRules {
		row := make(map[string]interface{})
		row["id"] = rt.ID
		row["script"] = rt.Rule
		routes[i] = row
		i++
	}

	return routes

}
