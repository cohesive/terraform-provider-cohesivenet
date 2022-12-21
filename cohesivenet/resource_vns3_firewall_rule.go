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
			"vns3": &schema.Schema{
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: getVns3AuthSchema(),
				},
			},
			"rule": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Nested Block for rules",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Id given to rule after it has been applied",
						},
						"script": &schema.Schema{
							Type:        schema.TypeString,
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
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

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
		resourceRulesDelete(ctx, d, m)
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	resourceRulesRead(ctx, d, m)

	return diags
}

func resourceRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	firewallResponse, err := c.GetFirewallRules()
	if err != nil {
		return diag.FromErr(err)
	}

	rules := flattenRulesData(firewallResponse)

	d.Set("rule", rules)

	return diags
}

func resourceRulesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	if d.HasChange("rule") {
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
		err := c.UpdateRules(ruleList)
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}
	return resourceRulesRead(ctx, d, m)
}

func resourceRulesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

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
