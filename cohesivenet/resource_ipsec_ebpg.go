package cohesivenet

import (
	"context"
	"strconv"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEbgp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEbgpCreate,
		ReadContext:   resourceEbgpRead,
		UpdateContext: resourceEbgpUpdate,
		DeleteContext: resourceEbgpDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"endpoint_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"ebgp_peer": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipaddress": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"asn": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"local_asn_alias": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"access_list": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"bgp_password": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"add_network_distance": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"add_network_distance_direction": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"add_network_distance_hops": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceEbgpCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	var diags diag.Diagnostics

	ebgp := d.Get("ebgp_peer").([]interface{})[0]
	bgp := ebgp.(map[string]interface{})

	ep := cn.EbgpPeer{

		Ipaddress:                   bgp["ipaddress"].(string),
		Asn:                         bgp["asn"].(int),
		LocalAsnAlias:               bgp["local_asn_alias"].(int),
		AccessList:                  bgp["access_list"].(string),
		AddNetworkDistanceHops:      bgp["add_network_distance_hops"].(int),
		BgpPassword:                 bgp["bgp_password"].(string),
		AddNetworkDistance:          bgp["add_network_distance"].(bool),
		AddNetworkDistanceDirection: bgp["add_network_distance_direction"].(string),
	}

	endId := d.Get("endpoint_id").(int)
	endpointId := strconv.Itoa(endId)

	_, err := c.CreateEbgpPeer(endpointId, &ep)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(endpointId)

	//resourceEndpointsRead(ctx, d, m)

	return diags
}

func resourceEbgpRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

/*
func resourceEbgpRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	var diags diag.Diagnostics

	endpointId := d.Id()

	endpoint, err := c.GetEndpoint(endpointId)
	if err != nil {
		return diag.FromErr(err)
	}

	flatEndpoint := flattenEndpointData(endpoint)

	if err := d.Set("endpoint", flatEndpoint); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(endpoint.Response.ID))

	return diags
}
*/
func resourceEbgpUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

/*
func resourceEbgpUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

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
*/
/*
func resourceEbgpDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
*/

func resourceEbgpDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	var diags diag.Diagnostics

	endpointId := d.Id()
	ebgpPeerId := d.Id()

	err := c.DeleteEbgpPeer(endpointId, ebgpPeerId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

/*
func flattenEbgpData(newEndpoint cn.NewEndpoint) []interface{} {
	endpoint := make([]interface{}, 1, 1)
	row := make(map[string]interface{})

	row["name"] = newEndpoint.Response.Name
	row["description"] = newEndpoint.Response.Description
	row["ipaddress"] = newEndpoint.Response.Ipaddress
	row["secret"] = newEndpoint.Response.Psk
	row["pfs"] = newEndpoint.Response.Pfs
	row["nat_t_enabled"] = newEndpoint.Response.NatTEnabled
	row["vpn_type"] = newEndpoint.Response.VpnType
	row["ike_version"] = newEndpoint.Response.IkeVersion
	row["route_based_int_address"] = newEndpoint.Response.RouteBasedIntAddress
	row["route_based_local"] = newEndpoint.Response.RouteBasedLocal
	row["route_based_remote"] = newEndpoint.Response.RouteBasedRemote
	row["extra_config"] = strings.Join(newEndpoint.Response.ExtraConfig, ", ")
	endpoint[0] = row

	return endpoint
}
*/
