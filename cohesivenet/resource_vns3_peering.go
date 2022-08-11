package cohesivenet

import (
	"context"
	"strconv"
	"time"
	"fmt"
	"strings"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet"
	macros "github.com/cohesive/cohesivenet-client-go/cohesivenet/macros"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)


func resourceVns3Peering() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePeeringCreate,
		ReadContext:   resourcePeeringRead,
		UpdateContext: resourceConfigUpdate,
		DeleteContext: resourceConfigDelete,
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
				ForceNew: true,
				Elem:     &schema.Resource{
					Schema: getVns3AuthSchema(),
				},
			},
			"peer": &schema.Schema{
				Type:     schema.TypeSet,
				MinItems: 1,
				Optional: true,
				Elem:     &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": &schema.Schema{
							Type:    schema.TypeString,
							Required: true,
						},
						"peer_id": &schema.Schema{
							Type:    schema.TypeInt,
							Required: true,
						},
						"overlay_mtu": &schema.Schema{
							Type:    schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func getPeeringStatus(ctx context.Context, vns3 *cn.VNS3Client) (*cn.PeersDetail, error) {
	apiRequest := vns3.PeeringApi.GetPeeringStatusRequest(ctx)
	resp, _, err := vns3.PeeringApi.GetPeeringStatus(apiRequest)
	if err != nil {
		return nil, err
	}
	peerDetail := resp.GetResponse()
	return &peerDetail, nil
}

func deleteAllPeers(ctx context.Context, vns3 *cn.VNS3Client) (*cn.PeersDetail, error) {
	deleteRequest := vns3.PeeringApi.DeletePeerRequest(ctx, 1)
	peerDetail, err := getPeeringStatus(ctx, vns3)
	resp, _, err := vns3.PeeringApi.DeletePeer(deleteRequest)

}

func resourcePeeringCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	peersSet, _ := d.Get("peer").(*schema.Set)
	for _, _peer := range peersSet.List() {
		peer := _peer.(map[string]any)
		peerName := peer["address"].(string)
		peerId := peer["peer_id"].(int32)
		peerRequest := cn.NewCreatePeerRequest(peerId, peerName)

		if peerMtu, hasMtu := peer["overlay_mtu"]; hasMtu {
			peerRequest.SetOverlayMtu(peerMtu.(string))
		}

		apiRequest := vns3.PeeringApi.PostCreatePeerRequest(ctx)
		apiRequest = apiRequest.CreatePeerRequest(*peerRequest)

		resp, _, err := vns3.PeeringApi.PostCreatePeer(apiRequest)
		if err != nil {

		}
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	resourcePeeringRead(ctx, d, m)
	return diags

	return diags
}

/*
func resourceConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

*/
func resourcePeeringRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}


	peerDetail, err := getPeering(ctx, vns3)

	configDetail, _, err := vns3.PeeringApi.GetConfig(vns3.ConfigurationApi.GetConfigRequest(ctx))
	if err != nil {
		return diag.FromErr(fmt.Errorf("VNS3 Config check error: %+v", err))
	}

	configData := configDetail.GetResponse()
	topologyChecksum := configData.GetTopologyChecksum()
	d.Set("topology_checksum", topologyChecksum)
	d.Set("licensed", configData.GetLicensed())

	keysetDetail, _, err := vns3.ConfigurationApi.GetKeyset(vns3.ConfigurationApi.GetKeysetRequest(ctx))

	if err != nil {
		return diag.FromErr(fmt.Errorf("VNS3 Keyset check error: %+v", err))
	}

	keysetData := keysetDetail.GetResponse()
	keysetChecksum := keysetData.GetChecksum()
	d.Set("keyset_checksum", keysetChecksum)

	return diags
}


func resourcePeeringUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// TODO: we could allow topology name and controller name to be reset and only fail 
	// when license params or keyset params change
	notsupportederror := fmt.Errorf("VNS3 config resource cannot be updated. Please redeploy a new server or reset defaults and edit terraform state")
	return diag.FromErr(notsupportederror)
}


func resourcePeeringDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	// Basically we just lie and say it was deleted.
	d.SetId("")

	return diags
}