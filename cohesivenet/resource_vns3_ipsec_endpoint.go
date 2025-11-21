package cohesivenet

import (
	"context"
	"strconv"
	"strings"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEndpoints() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEndpointsCreate,
		ReadContext:   resourceEndpointsRead,
		UpdateContext: resourceEndpointsUpdate,
		DeleteContext: resourceEndpointsDelete,
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
			"endpoint": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "Nested block for endpoint attributes",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name for new endpoint",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Description of new endpoint",
						},
						"ipaddress": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP address or remote device",
						},
						"secret": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Pre-shared key for IPSec connection",
						},
						"pfs": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Perfect Forward Secrecy setting. Default: false",
						},
						"ike_version": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "IKE version",
						},
						"nat_t_enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Perfect Forward Secrecy setting. Default: false",
						},
						"extra_config": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IPsec extra parameter settings for auth and encryption",
						},
						"private_ipaddress": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Remote Peer's IKE ID",
						},
						"vpn_type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type of VPN connection. VTI or GRE",
						},
						"route_based_int_address": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "If VTI a /30 address for the virtual interface",
						},
						"route_based_local": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Local subnet of IPsec tunnel",
						},
						"route_based_remote": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Remote subnet of IPsec tunnel",
						},
					},
				},
			},
		},
	}
}

func resourceEndpointsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	var diags diag.Diagnostics

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
		Extra_config:            strings.Replace(endpoint["extra_config"].(string), ",", " ", -1),
		Private_ipaddress:       endpoint["private_ipaddress"].(string),
		Vpn_type:                endpoint["vpn_type"].(string),
		Route_based_int_address: endpoint["route_based_int_address"].(string),
		Route_based_local:       endpoint["route_based_local"].(string),
		Route_based_remote:      endpoint["route_based_remote"].(string),
	}
	endpointResponse, err := c.CreateEndpoint(&ep)
	if err != nil {
		return diag.FromErr(err)
	}
	newEndpoint := endpointResponse.Response

	d.SetId(strconv.Itoa(newEndpoint.ID))

	resourceEndpointsRead(ctx, d, m)

	return diags
}

func resourceEndpointsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	var diags diag.Diagnostics

	endpointId := d.Id()

	endpoint, err := c.GetEndpoint(endpointId)
	if err != nil {
		return diag.FromErr(err)
	}

	flatEndpoint := flattenEndpointData(endpoint)

	if len(flatEndpoint) == 0 {
		d.SetId("")
		return diags
	}

	if err := d.Set("endpoint", flatEndpoint); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(endpoint.Response.ID))

	return diags
}

func resourceEndpointsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

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
			Extra_config:            strings.Replace(endpoint["extra_config"].(string), ",", " ", -1),
			Private_ipaddress:       endpoint["private_ipaddress"].(string),
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

func resourceEndpointsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	var diags diag.Diagnostics

	endpointId := d.Id()

	err := c.DeleteEndpoint(endpointId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func flattenEndpointData(newEndpoint cn.NewEndpoint) []interface{} {
	endpoint := make([]interface{}, 1)
	row := make(map[string]interface{})

	row["name"] = newEndpoint.Response.Name
	row["description"] = newEndpoint.Response.Description
	row["ipaddress"] = newEndpoint.Response.Ipaddress
	row["secret"] = newEndpoint.Response.Psk
	row["pfs"] = newEndpoint.Response.Pfs
	row["nat_t_enabled"] = newEndpoint.Response.NatTEnabled
	row["private_ipaddress"] = newEndpoint.Response.PrivateIpaddress
	row["vpn_type"] = newEndpoint.Response.VpnType
	row["ike_version"] = newEndpoint.Response.IkeVersion
	row["route_based_int_address"] = newEndpoint.Response.RouteBasedIntAddress
	row["route_based_local"] = newEndpoint.Response.RouteBasedLocal
	row["route_based_remote"] = newEndpoint.Response.RouteBasedRemote
	row["extra_config"] = strings.Join(newEndpoint.Response.ExtraConfig, ", ")
	endpoint[0] = row

	return endpoint
}
