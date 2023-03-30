package cohesivenet

import (
	"context"
	"fmt"
	"strconv"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVns3Peering() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePeeringCreate,
		ReadContext:   resourcePeeringRead,
		UpdateContext: resourcePeeringUpdate,
		DeleteContext: resourcePeeringDelete,
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
			"peer": &schema.Schema{
				Type:     schema.TypeSet,
				MinItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP address or DNS name of remote peer",
						},
						"peer_id": &schema.Schema{
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Peer id in topology",
						},
						"overlay_mtu": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "MTU overide for peering link",
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
	peerDetail, err := getPeeringStatus(ctx, vns3)
	if err != nil {
		return nil, err
	}

	for _, peer := range peerDetail.GetManagers() {
		if !peer.GetSelf() {
			deleteRequest := vns3.PeeringApi.DeletePeerRequest(ctx, peer.GetId())
			vns3.PeeringApi.DeletePeer(deleteRequest)
		}
	}

	return getPeeringStatus(ctx, vns3)
}

func createAllPeers(ctx context.Context, d *schema.ResourceData, vns3 *cn.VNS3Client) []string {
	peersSet, _ := d.Get("peer").(*schema.Set)
	failures := []string{}
	for _, _peer := range peersSet.List() {
		peer := _peer.(map[string]any)
		peerName := peer["address"].(string)
		peerId := int32(peer["peer_id"].(int))
		peerRequest := cn.NewCreatePeerRequest(peerId, peerName)
		peerMtu := peer["overlay_mtu"].(int)

		if peerMtu != 0 {
			vns3.Log.Info(fmt.Sprintf("Overlay MTU: %v", peerMtu))
			peerRequest.SetOverlayMtu(peerMtu)
		}

		apiRequest := vns3.PeeringApi.PostCreatePeerRequest(ctx).CreatePeerRequest(*peerRequest)
		_, _, err := vns3.PeeringApi.PostCreatePeer(apiRequest)
		if err != nil {
			apiError, isApiError := err.(*cn.GenericAPIError)
			// apiError := cn.ParseApiError(err)
			var errorString string
			if isApiError {
				errorMessage := apiError.GetErrorMessage()
				errorString = fmt.Sprintf("Create peer %v @ %v failed: %v", peerId, peerName, errorMessage)
			} else {
				errorString = fmt.Sprintf("Create peer %v @ %v failed: %v", peerId, peerName, err.Error())
			}
			failures = append(failures, errorString)
		}
	}

	return failures
}

func resourcePeeringCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	failures := createAllPeers(ctx, d, vns3)
	if len(failures) > 0 {
		errMessage := fmt.Sprintf("Failed to create all peers: %v.", failures)
		return diag.FromErr(fmt.Errorf(errMessage))
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	resourcePeeringRead(ctx, d, m)
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

	peerDetail, err := getPeeringStatus(ctx, vns3)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, peer := range peerDetail.GetManagers() {
		if !peer.GetSelf() {
			vns3.Log.Info(fmt.Sprintf("Read peer id=%v name=%v mtu=%v", peer.GetId(), peer.GetAddress(), peer.GetMtu()))
		}
	}

	return diags
}

func resourcePeeringUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	if d.HasChange("peer") {
		deleteAllPeers(ctx, vns3)
		failures := createAllPeers(ctx, d, vns3)
		if len(failures) != 0 {
			errMessage := fmt.Sprintf("Failed to create all peers: %v.", failures)
			return diag.FromErr(fmt.Errorf(errMessage))
		}
	}

	return diags
}

func resourcePeeringDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	deleteAllPeers(ctx, vns3)

	return diags
}
