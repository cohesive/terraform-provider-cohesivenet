package cohesivenet

import (
	"context"
	"strconv"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet/v1"
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
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CIDR of route",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of created route",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description of route",
						},
						"advertise": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Flag to advertise route to VNS3 Overlay Network",
						},
						"metric": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Route metric",
						},
						"enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Flag to enable route",
						},
						"netmask": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Netmask",
						},
						"editable": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Editable flag",
						},
						"table": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Table route is created in",
						},
						"interface": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies interface route applies to",
						},
					},
				},
			},
		},
	}
}

func dataSourceRoutesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}
	var diags diag.Diagnostics

	routeResponse, err := c.GetRoutes()
	if err != nil {
		return diag.FromErr(err)
	}

	routes := flattenRoutes(routeResponse)

	if err := d.Set("response", routes); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenRoutes(routeResponse cn.RouteResponse) interface{} {
	routes := make([]interface{}, len(routeResponse.Routes))

	i := 0
	for _, rt := range routeResponse.Routes {
		row := make(map[string]interface{})

		row["cidr"] = rt.Cidr
		row["id"] = rt.ID
		row["description"] = rt.Description
		row["advertise"] = rt.Advertise
		row["metric"] = rt.Metric
		row["enabled"] = rt.Enabled
		row["netmask"] = rt.Netmask
		row["editable"] = rt.Editable
		row["table"] = rt.Table
		row["interface"] = rt.Interface

		routes[i] = row
		i++
	}

	return routes

}
