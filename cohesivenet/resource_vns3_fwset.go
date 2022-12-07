package cohesivenet

import (
	"context"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFwSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFwSetCreate,
		ReadContext:   resourceFwSetRead,
		UpdateContext: resourceFwSetUpdate,
		DeleteContext: resourceFwSetDelete,
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
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of fwset",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Type of fwset",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of fwset",
			},
			"entries": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Nested fwset entries",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entry": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "fwset entry",
						},
						"comment": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "fwset entry comment",
						},
						"entry_resolved": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "fwset entry resolved",
						},
						"last_resolved": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "fwset entry last resolved",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "when fwset entry created",
						},
					},
				},
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When fwset was created",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update for fwset",
			},
			"size": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Size of fwset",
			},
		},
	}
}

func resourceFwSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	name := d.Get("name").(string)
	fwset_type := d.Get("type").(string)
	newFWSet := cn.NewCreateFirewallFwsetRequest(name, fwset_type)

	description, hasDescription := d.Get("description").(string)
	if hasDescription {
		newFWSet.SetDescription(description)
	}

	all_entries := d.Get("entries").([]interface{})
	newFWSet.SetEntries(parseFwSetEntries(all_entries))

	apiRequest := vns3.FirewallApi.PostCreateFirewallFwsetRequest(ctx)
	apiRequest = apiRequest.CreateFirewallFwsetRequest(*newFWSet)
	result, _, err := vns3.FirewallApi.PostCreateFirewallFwset(apiRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	apiResponse := result.GetResponse()
	d.SetId(apiResponse.GetName())

	resourceFwSetRead(ctx, d, m)

	return diags
}

func resourceFwSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	name := d.Get("name").(string)
	fwSetRequest := vns3.FirewallApi.GetFirewallFwsetRequest(ctx, name)
	fwSet, _, err := vns3.FirewallApi.GetFirewallFwset(fwSetRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	fwSetResponse := fwSet.GetResponse()

	d.Set("name", fwSetResponse.Name)
	d.Set("type", fwSetResponse.Type)
	d.Set("description", fwSetResponse.Description)
	d.Set("entries", flattenFwSetEntries(fwSetResponse.Entries))
	d.Set("created_at", fwSetResponse.CreatedAt.String())
	d.Set("updated_at", fwSetResponse.CreatedAt.String())
	d.Set("size", fwSetResponse.Size)

	return diags
}

func resourceFwSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	hasChange := false
	updateFwset := cn.NewUpdateFirewallFwsetRequest()

	if d.HasChange("name") {
		hasChange = true
		name := d.Get("name").(string)
		updateFwset.SetName(name)
	}

	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		updateFwset.SetDescription(description)
	}

	if d.HasChange("entries") {
		hasChange = true
		entries := d.Get("entries").([]interface{})
		updateFwset.Entries = parseFwSetEntries(entries)
	}

	if hasChange {
		apiRequest := vns3.FirewallApi.PutUpdateFirewallFwsetRequest(ctx, d.Id())
		apiRequest = apiRequest.UpdateFirewallFwsetRequest(*updateFwset)
		_, _, err := vns3.FirewallApi.PutUpdateFirewallFwset(apiRequest)
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	resourceFwSetRead(ctx, d, m)

	return diags
}

func resourceFwSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	id := d.Id()
	ruleRequest := vns3.FirewallApi.DeleteFirewallFwsetRequest(ctx, id)
	_, _, err := vns3.FirewallApi.DeleteFirewallFwset(ruleRequest)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func parseFwSetEntries(entries []interface{}) []cn.CreateFirewallEntryRequest {
	req_entries := []cn.CreateFirewallEntryRequest{}
	for _, entry := range entries {
		row := cn.CreateFirewallEntryRequest{}
		values := entry.(map[string]interface{})
		row.Entry = values["entry"].(string)
		comment := values["comment"].(string)
		row.Comment = &comment
		req_entries = append(req_entries, row)
	}
	return req_entries
}

func flattenFwSetEntries(entries []cn.FirewallFwsetEntry) interface{} {
	fwSetEntries := make([]interface{}, len(entries))

	for i, entry := range entries {
		row := make(map[string]interface{})
		row["entry"] = entry.Entry
		row["comment"] = entry.Comment
		row["entry_resolved"] = entry.EntryResolved
		row["last_resolved"] = entry.EntryResolved
		row["created_at"] = entry.CreatedAt.String()
		fwSetEntries[i] = row
	}

	return fwSetEntries

}
