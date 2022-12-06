package cohesivenet

import (
	"context"
	"fmt"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFirewallRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFirewallRulesCreate,
		ReadContext:   resourceFirewallRulesRead,
		UpdateContext: resourceFirewallRulesUpdate,
		DeleteContext: resourceFirewallRulesDelete,
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "Firewall rule in VNS3 syntax",
			},
			"position": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Position of firewall rule",
			},
			"table": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Firewall rule table in VNS3",
			},
			"rule_resolved": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Firewall rule resolved table in VNS3 syntax",
			},
			"comment": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Firewall comment",
			},
			"last_resolved": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Firewall rule date last resolved",
			},
			"disabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the rule is disabled or not",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Firewall rule created date",
			},
			"groups": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Firewall rule groups",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceFirewallRulesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	rule := d.Get("rule").(string)
	newRule := cn.NewCreateFirewallRuleRequest(rule)

	position, hasPosition := d.Get("position").(int32)
	if hasPosition {
		newRule.Position = &position
	}

	comment, hasComment := d.Get("comment").(string)
	if hasComment {
		newRule.Comment = &comment
	}

	groups, hasGroups := d.Get("groups").([]string)
	if hasGroups {
		newRule.Groups = groups
	}

	disabled, hasDisabled := d.Get("disabled").(bool)
	if hasDisabled {
		newRule.Disabled = &disabled
	}

	apiRequest := vns3.FirewallApi.PostCreateFirewallRuleRequest(ctx)
	apiRequest = apiRequest.CreateFirewallRuleRequest(*newRule)
	result, _, err := vns3.FirewallApi.PostCreateFirewallRule(apiRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	apiResponse := result.FirewallRuleDetailResponse.GetResponse()
	d.SetId(*apiResponse.Id)

	resourceFirewallRulesRead(ctx, d, m)

	return diags
}

func resourceFirewallRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	ruleRequest := vns3.FirewallApi.GetFirewallRulesRequest(ctx)
	rules, _, err := vns3.FirewallApi.GetFirewallRules(ruleRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	rulesResponse := rules.GetResponse()

	id := d.Id()
	rule := findRule(rulesResponse, id)

	if rule == nil {
		return diag.FromErr(fmt.Errorf("could not find firewall rule %v", id))
	}

	d.Set("rule", rule.Rule)
	d.Set("position", rule.Position)
	d.Set("table", rule.Table)
	d.Set("comment", rule.Comment)
	d.Set("created_at", rule.CreatedAt.String())
	d.Set("disabled", rule.Disabled)
	d.Set("groups", rule.Groups)
	d.Set("last_resolved", rule.LastResolved)
	d.Set("rule_resolved", rule.RuleResolved)

	return diags
}

func resourceFirewallRulesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	hasChange := false
	updateRule := cn.NewUpdateFirewallRuleRequest()

	if d.HasChange("rule") {
		hasChange = true
		rule := d.Get("rule").(string)
		updateRule.SetRule(rule)
	}

	if d.HasChange("disabled") {
		hasChange = true
		disabled := d.Get("disabled").(bool)
		updateRule.SetDisabled(disabled)
	}

	if d.HasChange("comment") {
		hasChange = true
		comment := d.Get("comment").(string)
		updateRule.SetComment(comment)
	}

	if d.HasChange("groups") {
		hasChange = true
		groups := d.Get("groups").([]string)
		updateRule.SetGroups(groups)
	}

	if hasChange {
		apiRequest := vns3.FirewallApi.PutUpdateFirewallRuleRequest(ctx, d.Id())
		apiRequest = apiRequest.UpdateFirewallRuleRequest(*updateRule)
		_, _, err := vns3.FirewallApi.PutUpdateFirewallRule(apiRequest)
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	resourceFirewallRulesRead(ctx, d, m)

	return diags
}

func resourceFirewallRulesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	id := d.Id()
	ruleRequest := vns3.FirewallApi.DeleteFirewallRuleRequest(ctx, id)
	_, _, err := vns3.FirewallApi.DeleteFirewallRule(ruleRequest)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func findRule(ruleResponse []cn.FirewallRule, id string) *cn.FirewallRule {

	for _, rt := range ruleResponse {
		if *rt.Id == id {
			return &rt
		}
	}
	return nil
}
