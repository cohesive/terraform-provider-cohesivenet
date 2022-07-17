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
						"row": &schema.Schema{
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
				},
			},
		},
	}
}

func dataSourceRoutesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cn.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	routeResponse, err := c.GetRoutes()
	if err != nil {
		return diag.FromErr(err)
	}

	//resp := make([]interface{}, 1, 1)

	newRoutes := routeResponse.Response.(map[string]interface{})

	//routes := flattenRoutes(&routeResponse)
	routes := flattenRoutes(newRoutes)

	if err := d.Set("response", routes); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenRoutes(routeResponse map[string]interface{}) interface{} {
	if routeResponse != nil {
		//routes := make(map[string]interface{})
		//routes := make([]interface{}, 26, 26)
		routes := make([]interface{}, len(routeResponse), len(routeResponse))

		i := 0
		//for id, rt := range routeResponse["response"].(map[string]interface{}) {
		for _, rt := range routeResponse {
			row := make(map[string]interface{})
			//rt_data := rt.(map[string]interface{})

			row["cidr"] = rt
			//row["id"] = id
			//			row["description"] = rt_data["description"]
			//			row["advertise"] = rt_data["advertise"]
			//			row["metric"] = rt_data["metric"]
			//			row["enabled"] = rt_data["enabled"]
			//			row["netmask"] = rt_data["netmask"]
			//			row["editable"] = rt_data["editable"]
			//			row["table"] = rt_data["table"]
			//			row["interface"] = rt_data["interface"]

			routes[i] = row
			i++
		}

		routes[i] = routes
		return routes
	}

	//return make(map[string]interface{})
	return make([]interface{}, 0)
}
