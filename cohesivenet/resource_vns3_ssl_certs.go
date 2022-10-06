package cohesivenet

import (
	"context"
	"log"
	"strconv"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSslCerts() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSslCertsCreate,
		ReadContext:   resourceSslCertsRead,
		UpdateContext: resourceSslCertsUpdate,
		DeleteContext: resourceSslCertsDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cert_file": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"key_file": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceSslCertsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	var diags diag.Diagnostics

	cert_file := d.Get("cert_file").(string)
	key_file := d.Get("key_file").(string)

	response, err := c.UpdateCerts(cert_file, key_file)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println(response)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	resourceSslCertsRead(ctx, d, m)

	return diags

}

func resourceSslCertsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceSslCertsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceSslCertsRead(ctx, d, m)
}

/*
func resourceSslCertsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
*/

func resourceSslCertsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceRoutesRead(ctx, d, m)
}

/*
func resourceSslCertsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	var diags diag.Diagnostics

	err := c.DeleteRoute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
*/
func flattenSslCertsData(routeResponse cn.RouteResponse) interface{} {
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
