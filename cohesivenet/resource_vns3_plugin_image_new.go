package cohesivenet

import (
	"context"
	"strconv"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePluginImageNew() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePluginImageCreateNew,
		ReadContext:   resourcePluginImageReadNew,
		UpdateContext: resourcePluginImageUpdateNew,
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

func parseImageResponseId(plugin cn.Plugin) string {
	iId := plugin.GetId()
	imageId := int(iId)
	imageIdString := strconv.Itoa(imageId)

	return imageIdString
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
	imageIdString := parseImageResponseId(imageData)
	/*
		imageResponse, err := c.CreateImage(&im)
		if err != nil {
			return diag.FromErr(err)
		}

		if len(imageResponse.Images) != 0 {
			d.SetId(imageResponse.Images[0].ID)
			resourcePluginImageRead(ctx, d, m)
		}
	*/

	d.SetId(imageIdString)
	resourcePluginImageReadNew(ctx, d, m)

	return diags
}

/*
func resourcePluginImageCreateNew(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	imageName := d.Get("image_name").(string)
	imageUrl := d.Get("url").(string)

	newImage := cn.NewCreateContainerImageRequest(imageName)
	newImage.SetUrl(imageUrl)
	imageDescription, hasDescription := d.Get("description").(string)
	if hasDescription {
		newImage.SetDescription(imageDescription)
	}

	//apiRequest := vns3.NewInstallPluginRequest

	var diags diag.Diagnostics

	return diags
}
*/
func resourcePluginImageReadNew(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}

/*
func resourcePluginImageReadNew(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	var diags diag.Diagnostics

	imageId := d.Id()

	img := d.Get("image").([]interface{})[0]
	im := img.(map[string]interface{})
	url := im["url"].(string)
	imageResponse, err := c.GetImage(imageId)
	if err != nil {
		return diag.FromErr(err)
	}

	image := flattenPluginImageDataNew(imageResponse, url)

	if err := d.Set("image", image); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(imageId)

	return diags
}
*/
func resourcePluginImageUpdateNew(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}

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

/*
func resourcePluginImageDeleteNew(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	var diags diag.Diagnostics

	id := d.Id()
	err := c.DeleteImage(id)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return diags
}

func flattenPluginImageDataNew(imageResponse cn.ImageResponse, url string) interface{} {
	image := make([]interface{}, len(imageResponse.Images))

	for _, ir := range imageResponse.Images {
		row := make(map[string]interface{})

		row["id"] = ir.ID
		row["image_name"] = ir.ImageName
		row["tag_name"] = ir.TagName
		row["status"] = ir.Status
		row["status_msg"] = ir.StatusMsg
		row["import_id"] = ir.ImportID
		row["created"] = ir.Created
		row["description"] = ir.Description
		row["comment"] = ir.Comment
		row["import_uuid"] = ir.ImportUUID
		row["url"] = url

		image[0] = row

	}

	return image

}
*/
