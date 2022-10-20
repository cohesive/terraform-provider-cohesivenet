package cohesivenet

import (
	"context"
	"strconv"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLink() *schema.Resource {
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
			"id": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Id for interface and link",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name for Link",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of transit link",
			},
			"conf": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Link conf (wireguard or openvpn)",
			},
			"policies": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional policies to place at end of conf file",
			},
		},
	}
}

func resourceLinkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics


	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	linkId := int32(d.Get("id").(int))
	linkName := d.Get("name").(string)
	linkConf := d.Get("conf").(string)
	linkDescription, hasDescription := d.Get("description").(string)
	linkPolicies, hasPolicies := d.Get("policies").(string)

	err := c.CreateRoute(routeList)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	resourceLinkRead(ctx, d, m)

	return diags

}

func resourceLinkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}


	routesResponse, err := c.GetRoutes()
	if err != nil {
		return diag.FromErr(err)
	}

	flatRoutes := flattenRouteData(routesResponse)

	d.Set("route", flatRoutes)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

func resourceLinkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
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

		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
		d.Set("last_updated", time.Now().Format(time.RFC850))

	}

	return resourceRoutesRead(ctx, d, m)
}

func resourceLinkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	err := c.DeleteRoute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
