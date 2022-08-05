package cohesivenet

import (
	"context"
	"strconv"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet/v1"
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
							Type:     schema.TypeString,
							Computed: true,
						},
						"running": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceContainerNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	containerResponse, err := c.GetContainerNetwork()
	if err != nil {
		return diag.FromErr(err)
	}

	resp := make([]interface{}, 1, 1)
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
