package cohesivenet

import (
	"context"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//Image V1 API and client - To be deprecated
func resourcePluginImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePluginImageCreate,
		ReadContext:   resourcePluginImageRead,
		// Currently update is not supported due to complexity.
		//UpdateContext: resourcePluginImageUpdate,
		DeleteContext: resourcePluginImageDelete,
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
			"image": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Nested block for image attributes",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Id of deployed image",
						},
						"image_name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Name of deployed image",
						},
						"url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "URL of the image file to be imported",
						},
						"buildurl": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "URL of a dockerfile that will be used to build the image",
						},
						"localbuild": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Local build file to create new image",
						},
						"localimage": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Local image to tag",
						},
						"imagefile": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Upload image file",
						},
						"buildfile": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Upload docker file or zipped docker context directory",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Description of deployed image",
						},
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Current status of upload",
						},
						"status_msg": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Status response",
						},
						"import_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Initial import Id",
						},
						"created": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "State of image",
						},
						"tag_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Image Tag",
						},
						"comment": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Comment",
						},
						"import_uuid": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Uuid of imported image",
						},
					},
				},
			},
		},
	}
}

func resourcePluginImageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	var diags diag.Diagnostics

	img := d.Get("image").([]interface{})[0]
	image := img.(map[string]interface{})

	im := cn.PluginImage{
		Name:        image["image_name"].(string),
		URL:         image["url"].(string),
		Buildurl:    image["buildurl"].(string),
		Localbuild:  image["localbuild"].(string),
		Localimage:  image["localimage"].(string),
		Imagefile:   image["imagefile"].(string),
		Buildfile:   image["buildfile"].(string),
		Description: image["description"].(string),
	}

	imageResponse, err := c.CreateImage(&im)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(imageResponse.Images) != 0 {
		d.SetId(imageResponse.Images[0].ID)
		resourcePluginImageRead(ctx, d, m)
	}
	return diags
}

func resourcePluginImageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	var diags diag.Diagnostics

	imageId := d.Id()

	img := d.Get("image").([]interface{})[0]
	im := img.(map[string]interface{})
	url := im["url"].(string)
	imageResponse, err := c.GetImage(imageId)
	if err != nil {
		return diag.FromErr(err)
	}

	image := flattenPluginImageData(imageResponse, url)

	if err := d.Set("image", image); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(imageId)

	return diags
}

func resourcePluginImageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	var diags diag.Diagnostics

	id := d.Id()
	err := c.DeleteImage(id)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return diags
}

func flattenPluginImageData(imageResponse cn.ImageResponse, url string) interface{} {
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
