package cohesivenet

// import (
// 	"context"
// 	"strconv"
// 	"strings"
// 	"time"

// 	cn "github.com/cohesive/cohesivenet-client-go/v1"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// )

// func findRouteInList(routes []interface{}, cidr string, routeInterface *string, metric *int) *map[string]interface{} {

// 	for k, v := range kvs {
//         fmt.Printf("%s -> %s\n", k, v)
//     }
// }

// func resourceRoute() *schema.Resource {
// 	return &schema.Resource{
// 		CreateContext: resourceRouteCreate,
// 		ReadContext:   resourceRouteRead,
// 		UpdateContext: resourceRouteUpdate,
// 		DeleteContext: resourceRouteDelete,
// 		Schema: map[string]*schema.Schema{
// 			"last_updated": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"id": &schema.Schema{
// 				Type:     schema.TypeInt,
// 				Computed: true,
// 			},
// 			"netmask": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Computed: true,
// 			},
// 			"cidr": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Required: true,
// 				ForceNew: true,
// 			},
// 			"interface": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				ForceNew: true,
// 			},
// 			"gateway": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				ForceNew: true,
// 			},
// 			"description": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				ForceNew: true, // i wish our api supported updating description without destroy. or we could remove
// 			},
// 			// "enabled": &schema.Schema{
// 				// Type:     schema.TypeBool,
// 				// Optional: true,
// 				// Computed: true,
// 			// },
// 			"table": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Default: "main"
// 			},
// 			"metric": &schema.Schema{
// 				Type:     schema.TypeInt,
// 				Optional: true,
// 			},
// 			"advertise": &schema.Schema{
// 				Type:     schema.TypeBool,
// 				Optional: true,
// 				Default: true,
// 			},
// 			"tunnel": &schema.Schema{
// 				Type:     schema.TypeInt,
// 				Optional: true,
// 			},
// 		},
// 	}
// }

// func resourceRouteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	vns3 := m["vns3"].(cohesivenet.VNS3Client)
// 	auth := m["auth"].(context.Context)

// 	// Warning or errors can be collected in a slice type
// 	var diags diag.Diagnostics

// 	allRoutesReq := vns3.RoutingApi.GetRoutes(auth)
// 	allRoutes, _, getErr := vns3.RoutingApi.GetRoutesExecute(allRoutesReq)

// 	req := cohesivenet.NewCreateRouteRequest(d.Get("cidr"))

// 	if routeInt := d.Get("interface"); routeInt != nil {
// 		req.Interface = routeInt
// 	}

// 	if description := d.Get("description"); description != nil {
// 		req.Description = description
// 	}

// 	if gateway := d.Get("gateway"); gateway != nil {
// 		req.Gateway = gateway
// 	}

// 	if table := d.Get("table"); table != nil {
// 		req.Table = table
// 	}

// 	if tunnelid := d.Get("tunnel"); tunnelid != nil {
// 		req.Tunnel = tunnelid
// 	}

// 	if metric := d.Get("metric"); metric != nil {
// 		req.Metric = metric
// 	}
	
// 	if advertiseRoute := d.Get("advertise"); advertiseRoute != nil {
// 		req.Advertise = advertiseRoute
// 	}

// 	apiReq := vns3.RoutingApi.PostCreateRoute(auth)
// 	apiReq.createRouteRequest = req
// 	resp, _, err := vns3.RoutingApi.PostCreateRouteExecute(apiReq)

// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
// 	newEndpoint := endpointResponse.Response

// 	d.SetId(strconv.Itoa(newEndpoint.ID))

// 	resourceRouteRead(ctx, d, m)

// 	return diags
// }

// /*
// func resourceRouteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	// Warning or errors can be collected in a slice type
// 	var diags diag.Diagnostics

// 	return diags
// }

// */
// func resourceRouteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	c := m.(*cn.Client)
// 	// Warning or errors can be collected in a slice type
// 	var diags diag.Diagnostics

// 	endpointId := d.Id()

// 	endpoint, err := c.GetEndpoint(endpointId)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}

// 	//newEndpoint := endpoints.Response.(map[string]interface{})
// 	flatEndpoint := flattenEndpointData(endpoint)

// 	if err := d.Set("endpoint", flatEndpoint); err != nil {
// 		return diag.FromErr(err)
// 	}

// 	// always run
// 	//d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
// 	d.SetId(strconv.Itoa(endpoint.Response.ID))

// 	return diags
// }

// /*
// func resourceRouteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	return resourceRouteRead(ctx, d, m)
// }
// */

// func resourceRouteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	c := m.(*cn.Client)

// 	endpointId := d.Id()

// 	if d.HasChange("endpoint") {

// 		endp := d.Get("endpoint").([]interface{})[0]
// 		endpoint := endp.(map[string]interface{})

// 		ep := cn.Endpoint{
// 			Name:                    endpoint["name"].(string),
// 			Description:             endpoint["description"].(string),
// 			Ipaddress:               endpoint["ipaddress"].(string),
// 			Secret:                  endpoint["secret"].(string),
// 			Pfs:                     endpoint["pfs"].(bool),
// 			Ike_version:             endpoint["ike_version"].(int),
// 			Nat_t_enabled:           endpoint["nat_t_enabled"].(bool),
// 			Extra_config:            endpoint["extra_config"].(string),
// 			Vpn_type:                endpoint["vpn_type"].(string),
// 			Route_based_int_address: endpoint["route_based_int_address"].(string),
// 			Route_based_local:       endpoint["route_based_local"].(string),
// 			Route_based_remote:      endpoint["route_based_remote"].(string),
// 		}

// 		_, err := c.UpdateEndpoint(endpointId, &ep)
// 		if err != nil {
// 			return diag.FromErr(err)
// 		}

// 		d.Set("last_updated", time.Now().Format(time.RFC850))
// 	}

// 	return resourceRouteRead(ctx, d, m)
// }

// func resourceRouteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	c := m.(*cn.Client)

// 	// Warning or errors can be collected in a slice type
// 	var diags diag.Diagnostics

// 	endpointId := d.Id()

// 	err := c.DeleteEndpoint(endpointId)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}

// 	d.SetId("")

// 	return diags
// }

// func flattenEndpointData(newEndpoint cn.NewEndpoint) []interface{} {
// 	endpoint := make([]interface{}, 1, 1)
// 	row := make(map[string]interface{})

// 	row["name"] = newEndpoint.Response.Name
// 	row["description"] = newEndpoint.Response.Description
// 	row["ipaddress"] = newEndpoint.Response.Ipaddress
// 	row["secret"] = newEndpoint.Response.Psk
// 	row["pfs"] = newEndpoint.Response.Pfs
// 	row["nat_t_enabled"] = newEndpoint.Response.NatTEnabled
// 	row["vpn_type"] = newEndpoint.Response.VpnType
// 	row["ike_version"] = newEndpoint.Response.IkeVersion
// 	row["route_based_int_address"] = newEndpoint.Response.RouteBasedIntAddress
// 	row["route_based_local"] = newEndpoint.Response.RouteBasedLocal
// 	row["route_based_remote"] = newEndpoint.Response.RouteBasedRemote
// 	row["extra_config"] = strings.Join(newEndpoint.Response.ExtraConfig, ", ")
// 	endpoint[0] = row

// 	return endpoint
// }
