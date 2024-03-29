package cohesivenet

import (
	"context"
	"strconv"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet/v1"
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
			"vns3": &schema.Schema{
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: getVns3AuthSchema(),
				},
			},
			"route": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "Nested block for route attributes",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "CIDR of route",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
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
							Optional:    true,
							Default:     false,
							Description: "Flag to advertise route to VNS3 Overlay Network",
						},
						"metric": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Route metric",
						},
						"enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
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
							Optional:    true,
							Computed:    true,
							Description: "Editable flag",
						},
						"table": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Table route is created in",
						},
						"interface": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies interface route applies to",
						},
						"gateway": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies network gateway",
						},
						"tunnel": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "GRE endpoint id (if applicable)",
						},
					},
				},
			},
		},
	}
}

func resourceRoutesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	var routeList []*cn.Route
	routes := d.Get("route").([]interface{})
	for _, route := range routes {
		rt := route.(map[string]interface{})
		route := cn.Route{
			Cidr:        rt["cidr"].(string),
			Description: rt["description"].(string),
			Interface:   rt["interface"].(string),
			Gateway:     rt["gateway"].(string),
			Tunnel:      rt["tunnel"].(int),
			Advertise:   rt["advertise"].(bool),
			Metric:      rt["metric"].(int),
		}

		routeList = append(routeList, &route)
	}
	err := c.CreateRoute(routeList)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	resourceRoutesRead(ctx, d, m)
	return diags

}

func resourceRoutesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	routesResponse, err := c.GetRoutes()
	if err != nil {
		return diag.FromErr(err)
	}

	flatRoutes := flattenRouteData(routesResponse)
	d.Set("route", flatRoutes)
	return diags
}

func resourceRoutesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	if d.HasChange("route") {
		var routeList []*cn.Route
		routes := d.Get("route").([]interface{})
		for _, route := range routes {
			rt := route.(map[string]interface{})
			route := cn.Route{
				Cidr:        rt["cidr"].(string),
				Description: rt["description"].(string),
				Interface:   rt["interface"].(string),
				Gateway:     rt["gateway"].(string),
				Tunnel:      rt["tunnel"].(int),
				Advertise:   rt["advertise"].(bool),
				Metric:      rt["metric"].(int),
			}

			routeList = append(routeList, &route)
		}
		err := c.UpdateRoute(routeList)
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))

	}
	return resourceRoutesRead(ctx, d, m)
}

func resourceRoutesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}
	err := c.DeleteRoute()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func flattenRouteData(routeResponse cn.RouteResponse) interface{} {
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
		row["gateway"] = rt.Gateway
		row["tunnel"] = rt.Tunnel
		row["interface"] = rt.Interface
		routes[i] = row
		i++
	}
	return routes
}
