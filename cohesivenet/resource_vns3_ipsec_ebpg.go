package cohesivenet

import (
	"context"
	"strconv"
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
				Optional: true,
				Computed: true,
			},
			"ebgp_peer": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ipaddress": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"asn": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"local_asn_alias": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"access_list": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"bgp_password": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"add_network_distance": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"add_network_distance_direction": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"add_network_distance_hops": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
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

/*
func resourceEbgpRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
*/

func resourceEbgpRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	var diags diag.Diagnostics

	//endpointId := d.Id()
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

/*
func resourceEbgpUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
*/

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
			AccessList:                  bgp["access_list"].(string),
			AddNetworkDistanceHops:      bgp["add_network_distance_hops"].(int),
			BgpPassword:                 bgp["bgp_password"].(string),
			AddNetworkDistance:          bgp["add_network_distance"].(bool),
			AddNetworkDistanceDirection: bgp["add_network_distance_direction"].(string),
		}
		_, err := c.UpdateEbgpPeer(endpointId, ebgpPeerId, &ep)
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceEbgpRead(ctx, d, m)
}

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

func flattenEbgpData(newPeer cn.EbgpPeer) []interface{} {
	epeer := make([]interface{}, 1, 1)
	row := make(map[string]interface{})

	row["ipaddress"] = newPeer.Ipaddress
	row["asn"] = newPeer.Asn
	row["local_asn_alias"] = newPeer.LocalAsnAlias
	row["access_list"] = newPeer.AccessList
	row["add_network_distance_hops"] = newPeer.AddNetworkDistanceHops
	row["bgp_password"] = newPeer.BgpPassword
	//row["bgp_password"] = "testtesticle"
	row["add_network_distance"] = newPeer.AddNetworkDistance
	row["add_network_distance_direction"] = newPeer.AddNetworkDistanceDirection
	epeer[0] = row

	return epeer
}