package cohesivenet

import (
	"context"
	"log"
	"strconv"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//Plugin image V2 API and go client
func resourcePluginImageNew() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePluginImageCreateNew,
		ReadContext:   resourcePluginImageReadNew,
		//Update TODO
		//UpdateContext: resourcePluginImageUpdateNew,
		DeleteContext: resourcePluginImageDeleteNew,
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
				Elem: &schema.Resource{
					Schema: getVns3AuthSchema(),
				},
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of deployed image",
			},
			"image_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "URL of the image file to be imported",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Description of deployed image",
			},
			"command": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "URL of a dockerfile that will be used to build the image",
			},
			"documentation_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Local build file to create new image",
			},
			"support_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Local image to tag",
			},
			"catalog_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Upload image file",
			},
			"version": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Upload docker file or zipped docker context directory",
			},
		},
	}
}

func parsePluginResponseId(plugin cn.Plugin) (string, int32) {
	iId := plugin.GetId()
	imageId := int(iId)
	imageIdString := strconv.Itoa(imageId)

	return imageIdString, iId
}

func resourcePluginImageCreateNew(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	name := d.Get("name").(string)
	image_url := d.Get("image_url").(string)
	description := d.Get("description").(string)
	command := d.Get("command").(string)

	newImage := cn.NewInstallPluginRequest(name, image_url)
	newImage.SetDescription(description)
	newImage.SetCommand(command)

	apiRequest := vns3.NetworkEdgePluginsApi.InstallPluginRequest(ctx)
	apiRequest = apiRequest.InstallPluginRequest(*newImage)
	detail, _, err := vns3.NetworkEdgePluginsApi.InstallPlugin(apiRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	imageData := detail.GetResponse()
	imageIdString, iId := parsePluginResponseId(imageData)

	timer := time.Tick(5 * time.Second)
	for _ = range timer {
		imageStatus := vns3.NetworkEdgePluginsApi.GetPluginRequest(ctx, iId)
		imageDetail, _, err := vns3.NetworkEdgePluginsApi.GetPlugin(imageStatus)
		if err != nil {
			return diag.FromErr(err)
		}
		plugin := imageDetail.GetResponse()
		pluginStatus := plugin.GetStatus()
		if pluginStatus == "ready" {
			log.Println(pluginStatus)
			d.SetId(imageIdString)
			resourcePluginImageReadNew(ctx, d, m)
			break
		}

	}

	return diags
}

func resourcePluginImageReadNew(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	Id := d.Id()
	iId, _ := strconv.Atoi(Id)
	PluginId := int32(iId)

	imageStatus := vns3.NetworkEdgePluginsApi.GetPluginRequest(ctx, PluginId)
	imageDetail, _, err := vns3.NetworkEdgePluginsApi.GetPlugin(imageStatus)
	if err != nil {
		return diag.FromErr(err)
	}
	plugin := imageDetail.GetResponse()
	pluginIdString, _ := parsePluginResponseId(plugin)
	d.SetId(pluginIdString)

	return diags
}

/* TODO - Not sure of value TBD
func resourcePluginImageUpdateNew(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}
*/
func resourcePluginImageDeleteNew(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	Id := d.Id()
	iId, _ := strconv.Atoi(Id)
	imageId := int32(iId)

	apiRequest := vns3.NetworkEdgePluginsApi.DeletePluginRequest(ctx, imageId)
	_, _, err := vns3.NetworkEdgePluginsApi.DeletePlugin(apiRequest)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
