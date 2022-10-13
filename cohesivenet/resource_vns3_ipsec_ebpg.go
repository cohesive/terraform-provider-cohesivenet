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
				Computed: true,
			},
			"ebgp_peer": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "Nested block for eBGP peer attributes",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the eBGP peer",
						},
						"ipaddress": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP address or neighbor IP for BGP",
						},
						"asn": &schema.Schema{
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Autonomous System Number of your network",
						},
						"local_asn_alias": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "ASN alias",
						},
						"access_list": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Access Control List. IN PERMIT xxxx / OUT PERMIT xxxx",
						},
						"bgp_password": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Password for BGP, if required",
						},
						"add_network_distance": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Specifies if we are using network distance weighting, Default: false",
						},
						"add_network_distance_direction": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies direction for distance weighting. IN / OUT",
						},
						"add_network_distance_hops": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies how many hops for network distance weighting",
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
		AccessList:                  strings.Replace(bgp["access_list"].(string), ",", "\n", -1),
		AddNetworkDistanceHops:      bgp["add_network_distance_hops"].(int),
		BgpPassword:                 bgp["bgp_password"].(string),
		AddNetworkDistance:          bgp["add_network_distance"].(bool),
		AddNetworkDistanceDirection: bgp["add_network_distance_direction"].(string),
	}

	endId := d.Get("endpoint_id").(int)
	endpointId := strconv.Itoa(endId)

	peerResponse, err := c.CreateEbgpPeer(endpointId, &ep)
	if err != nil {
		return diag.FromErr(err)
	}

	peerId := peerResponse.ID
	d.SetId(peerId)
	//d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	resourceEbgpRead(ctx, d, m)

	return diags
}

func resourceEbgpRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	var diags diag.Diagnostics

	endId := d.Get("endpoint_id").(int)
	endpointId := strconv.Itoa(endId)

	ebgpPeerId := d.Id()

	ebgPeer, err := c.GetEbgpPeer(endpointId, ebgpPeerId)
	if err != nil {
		return diag.FromErr(err)
	}

	flatPeer := flattenEbgpData(ebgPeer)

	if err := d.Set("ebgp_peer", flatPeer); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ebgpPeerId)
	//d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func resourceEbgpUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	endId := d.Get("endpoint_id").(int)
	endpointId := strconv.Itoa(endId)
	ebgpPeerId := d.Id()

	if d.HasChange("ebgp_peer") {

		ebgp := d.Get("ebgp_peer").([]interface{})[0]
		bgp := ebgp.(map[string]interface{})

		ep := cn.EbgpPeer{
			Ipaddress:                   bgp["ipaddress"].(string),
			Asn:                         bgp["asn"].(int),
			LocalAsnAlias:               bgp["local_asn_alias"].(int),
			AccessList:                  strings.Replace(bgp["access_list"].(string), ",", "\n", -1),
			AddNetworkDistanceHops:      bgp["add_network_distance_hops"].(int),
			BgpPassword:                 bgp["bgp_password"].(string),
			AddNetworkDistance:          bgp["add_network_distance"].(bool),
			AddNetworkDistanceDirection: bgp["add_network_distance_direction"].(string),
		}
		newPeer, err := c.UpdateEbgpPeer(endpointId, ebgpPeerId, &ep)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("id", newPeer.ID)
		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceEbgpRead(ctx, d, m)
}

func resourceEbgpDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	var diags diag.Diagnostics

	endId := d.Get("endpoint_id").(int)
	ebgpPeerId := d.Id()
	endpointId := strconv.Itoa(endId)
	err := c.DeleteEbgpPeer(endpointId, ebgpPeerId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func flattenEbgpData(newPeer cn.EbgpPeer) []interface{} {
	epeer := make([]interface{}, 1, 1)
	row := make(map[string]interface{})

	row["ipaddress"] = newPeer.Ipaddress
	row["asn"] = newPeer.Asn
	row["local_asn_alias"] = newPeer.LocalAsnAlias
	row["access_list"] = newPeer.AccessList
	row["add_network_distance_hops"] = newPeer.AddNetworkDistanceHops
	row["bgp_password"] = newPeer.BgpPassword
	row["add_network_distance"] = newPeer.AddNetworkDistance
	row["add_network_distance_direction"] = newPeer.AddNetworkDistanceDirection
	epeer[0] = row

	return epeer
}
