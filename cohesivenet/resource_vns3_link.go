package cohesivenet

import (
	"context"
	"time"
	"strings"
	"fmt"

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
				ForceNew:    true,
				Description: "Link conf (wireguard or openvpn)",
			},
			"policies": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional policies to place at end of conf file",
			},
			"clientpack_ip": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Overlay IP address for link",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of client connection",
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

    newLink := cn.NewCreateLinkRequest(linkId, linkName)
    // confBytes, err := os.ReadFile("openvpn.conf")
    // if err != nil {
    //     // fmt.Print(err)
    //     os.Exit(1)
    // }

    newLink.SetConf(linkConf)
	linkDescription, hasDescription := d.Get("description").(string)
	if hasDescription {
		newLink.SetDescription(linkDescription)
	}

	linkPolicies, hasPolicies := d.Get("policies").(string)
	if hasPolicies {
		policiesList := strings.Split(linkPolicies, "\n")
		newLink.SetPolicies(policiesList)
	}

	apiRequest := vns3.OverlayNetworkApi.CreateLinkRequest(ctx)
	apiRequest = apiRequest.CreateLinkRequest(*newLink)
	detail, _, err := vns3.OverlayNetworkApi.CreateLink(apiRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	linkData := detail.GetResponse()

	d.SetId(string(linkData.GetId()))

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

	linkId := d.Id()
	vns3.Log.Info(fmt.Sprintf("Reading linkId %v", string(linkId)))
	// if linkId != "" {

	// }

	getLinkRequest := vns3.OverlayNetworkApi.GetLinkRequest(ctx, linkId)
	detail, _, err := vns3.OverlayNetworkApi.GetLink(getLinkRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("VNS3 GET Link error: %+v", err))
	}

	link := detail.GetResponse()
	d.Set("clientpack_ip", link.GetClientpackIp())
	d.Set("type", link.GetType())
	d.SetId(string(link.GetId()))
	return diags
}

func resourceLinkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	hasChange := false
	updateBody := cn.NewUpdateLinkRequest()

	if d.HasChange("policies") {
		hasChange = true
		linkPolicies := d.Get("policies").(string)
		policiesList := strings.Split(linkPolicies, "\n")
		updateBody.SetPolicies(policiesList)
	}

	if d.HasChange("name") {
		hasChange = true
		name := d.Get("name").(string)
		updateBody.SetName(name)
	}

	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		updateBody.SetDescription(description)
	}

	if hasChange {
		linkId := d.Id()
		apiRequest := vns3.OverlayNetworkApi.PutUpdateLinkRequest(ctx, linkId)
		apiRequest = apiRequest.UpdateLinkRequest(*updateBody)
		_, _, err := vns3.OverlayNetworkApi.PutUpdateLink(apiRequest)
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	resourceLinkRead(ctx, d, m)
	return diags
}

func resourceLinkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	linkId := d.Id()
	apiRequest := vns3.OverlayNetworkApi.DeleteLinkRequest(ctx, linkId)
	_, _, err := vns3.OverlayNetworkApi.DeleteLink(apiRequest)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
