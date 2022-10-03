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
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	var diags diag.Diagnostics
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
	_, err := c.CreateRoute(routeList)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	resourceRoutesRead(ctx, d, m)

	return diags

}

func resourceRoutesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	var diags diag.Diagnostics

	//routeId := d.Id()

	routesResponse, err := c.GetRoutes()
	if err != nil {
		return diag.FromErr(err)
	}

	flatRoutes := flattenRouteData(routesResponse)

	d.Set("route", flatRoutes)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

/*
func resourceRoutesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceRoutesRead(ctx, d, m)
}
*/
func resourceRoutesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

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
		_, err := c.UpdateRoute(routeList)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
		d.Set("last_updated", time.Now().Format(time.RFC850))

	}

	return resourceRoutesRead(ctx, d, m)
}

func resourceRoutesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	var diags diag.Diagnostics

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
		row["interface"] = rt.Interface

		routes[i] = row
		i++
	}

	return routes

}
