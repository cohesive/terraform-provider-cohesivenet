package cohesivenet

import (
	"context"
	"strconv"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRoutes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRoutesRead,
		Schema: map[string]*schema.Schema{
			"response": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"advertise": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"metric": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"enabled": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"netmask": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"editable": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"table": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"interface": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRoutesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cn.Client)

	var diags diag.Diagnostics

	routeResponse, err := c.GetRoutes()
	if err != nil {
		return diag.FromErr(err)
	}

	newRoutes := routeResponse.Response.(map[string]interface{})

	routes := flattenRoutes(newRoutes)

	if err := d.Set("response", routes); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenRoutes(routeResponse map[string]interface{}) interface{} {
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
