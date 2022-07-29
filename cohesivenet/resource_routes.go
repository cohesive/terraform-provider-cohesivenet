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

	// Warning or errors can be collected in a slice type
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
	//idx := 0
	for _, r := range routes {
		values := r.(map[string]interface{})
		value, _ := strconv.Atoi(values["id"].(string))
		if value > highest {
			//idx = count
			highest = value
		}
	}

	//rId := routes[idx].(map[string]interface{})["id"]

	d.SetId(strconv.Itoa(highest))

	//d.SetId(strconv.Itoa(routeId["id"].(string)))
	//d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	//resourceEndpointsRead(ctx, d, m)

	return diags
}

func resourceRoutesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

/*
func resourceRoutesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cn.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	endpointId := d.Id()

	endpoint, err := c.GetEndpoint(endpointId)
	if err != nil {
		return diag.FromErr(err)
	}

	//newEndpoint := endpoints.Response.(map[string]interface{})
	flatEndpoint := flattenEndpointData(endpoint)

	if err := d.Set("endpoint", flatEndpoint); err != nil {
		return diag.FromErr(err)
	}

	// always run
	//d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.SetId(strconv.Itoa(endpoint.Response.ID))

	return diags
}
*/

func resourceRoutesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceEndpointsRead(ctx, d, m)
}

/*
func resourceRoutesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cn.Client)

	endpointId := d.Id()

	if d.HasChange("endpoint") {

		endp := d.Get("endpoint").([]interface{})[0]
		endpoint := endp.(map[string]interface{})

		ep := cn.Endpoint{
			Name:                    endpoint["name"].(string),
			Description:             endpoint["description"].(string),
			Ipaddress:               endpoint["ipaddress"].(string),
			Secret:                  endpoint["secret"].(string),
			Pfs:                     endpoint["pfs"].(bool),
			Ike_version:             endpoint["ike_version"].(int),
			Nat_t_enabled:           endpoint["nat_t_enabled"].(bool),
			Extra_config:            endpoint["extra_config"].(string),
			Vpn_type:                endpoint["vpn_type"].(string),
			Route_based_int_address: endpoint["route_based_int_address"].(string),
			Route_based_local:       endpoint["route_based_local"].(string),
			Route_based_remote:      endpoint["route_based_remote"].(string),
		}

		_, err := c.UpdateEndpoint(endpointId, &ep)
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceEndpointsRead(ctx, d, m)
}

func resourceRoutesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
*/

func resourceRoutesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cn.Client)

	// Warning or errors can be collected in a slice type
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
