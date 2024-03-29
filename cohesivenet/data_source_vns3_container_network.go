package cohesivenet

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceContainerNetwork() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerNetworkRead,
		Schema: map[string]*schema.Schema{
			"response": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VNS3 plugin network subnet",
						},
						"running": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "VNS3 plugun network state",
						},
					},
				},
			},
		},
	}
}

func dataSourceContainerNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	var diags diag.Diagnostics

	containerResponse, err := c.GetContainerNetwork()
	if err != nil {
		return diag.FromErr(err)
	}

	resp := make([]interface{}, 1)
	row := make(map[string]interface{})

	row["network"] = containerResponse.Response.Network
	row["running"] = containerResponse.Response.Running
	resp[0] = row

	if err := d.Set("response", resp); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
