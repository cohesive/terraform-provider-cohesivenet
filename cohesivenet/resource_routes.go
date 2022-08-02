package cohesivenet

import (
	"context"
	"strconv"

	cn "github.com/cohesive/cohesivenet-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRoutes() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoutesCreate,
		ReadContext:   resourceRoutesRead,
		UpdateContext: resourceRoutesUpdate,
		DeleteContext: resourceRoutesDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"route": &schema.Schema{
				Type:     schema.TypeList,
				ForceNew: true,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"id": &schema.Schema{
							Type:     schema.TypeInt,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"advertise": &schema.Schema{
							Type:     schema.TypeBool,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"metric": &schema.Schema{
							Type:     schema.TypeInt,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"enabled": &schema.Schema{
							Type:     schema.TypeBool,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"netmask": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"editable": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"table": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"interface": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"gateway": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"tunnel": &schema.Schema{
							Type:     schema.TypeInt,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceRoutesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cn.Client)

	var diags diag.Diagnostics

	rte := d.Get("route").([]interface{})[0]
	route := rte.(map[string]interface{})

	rt := cn.Route{
		Cidr:        route["cidr"].(string),
		Description: route["description"].(string),
		Interface:   route["interface"].(string),
		Gateway:     route["gateway"].(string),
		Tunnel:      route["tunnel"].(int),
		Advertise:   route["advertise"].(bool),
		Metric:      route["metric"].(int),
	}

	routeResponse, err := c.CreateRoute(&rt)
	if err != nil {
		return diag.FromErr(err)
	}

	routes := flattenRouteData(routeResponse)

	highest := 0
	for _, r := range routes {
		values := r.(map[string]interface{})
		value, _ := strconv.Atoi(values["id"].(string))
		if value > highest {
			highest = value
		}
	}

	d.SetId(strconv.Itoa(highest))

	resourceRoutesRead(ctx, d, m)

	return diags
}

func resourceRoutesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cn.Client)

	var diags diag.Diagnostics

	routeId := d.Id()

	routeResponse, err := c.GetRoute(routeId)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(routeResponse.Routes[0].ID)
	return diags
}

func resourceRoutesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceEndpointsRead(ctx, d, m)
}

func resourceRoutesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cn.Client)

	var diags diag.Diagnostics

	routeId := d.Id()

	err := c.DeleteRoute(routeId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func flattenRouteData(routeResponse cn.RouteResponse) []interface{} {
	routes := make([]interface{}, len(routeResponse.Routes), len(routeResponse.Routes))

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
