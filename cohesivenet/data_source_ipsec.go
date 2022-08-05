package cohesivenet

import (
	"context"
	"strconv"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEndpoints() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEndpointsRead,
		Schema: map[string]*schema.Schema{
			"response": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"local_subnet": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_subnet": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpointid": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"endpoint_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"active": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enabled": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"connected": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceEndpointsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	endpointResponse, err := c.GetEndpoints()
	if err != nil {
		return diag.FromErr(err)
	}

	newEndpoints := endpointResponse.Endpoints.(map[string]interface{})

	endpoints := flattenEndpoints(newEndpoints)

	if err := d.Set("response", endpoints); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenEndpoints(endpointResponse map[string]interface{}) interface{} {
	if endpointResponse != nil {
		endpoints := make([]interface{}, len(endpointResponse), len(endpointResponse))

		i := 0
		for _, rt := range endpointResponse {
			row := make(map[string]interface{})
			ep_data := rt.(map[string]interface{})

			row["local_subnet"] = ep_data["local_subnet"]
			row["remote_subnet"] = ep_data["remote_subnet"]
			row["endpointid"] = ep_data["endpointid"]
			row["endpoint_name"] = ep_data["endpoint_name"]
			row["active"] = ep_data["active"].(bool)
			row["enabled"] = ep_data["enabled"].(bool)
			row["connected"] = ep_data["connected"].(bool)

			endpoints[i] = row
			i++
		}

		return endpoints
	}

	return make([]interface{}, 0)
}
