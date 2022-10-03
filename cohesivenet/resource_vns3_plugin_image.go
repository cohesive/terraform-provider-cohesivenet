package cohesivenet

import (
	"context"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePluginImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePluginImageCreate,
		ReadContext:   resourcePluginImageRead,
		UpdateContext: resourcePluginImageUpdate,
		DeleteContext: resourcePluginImageDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"image": &schema.Schema{
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
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"url": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"buildurl": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"localbuild": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"localimage": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"imagefile": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"buildfile": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"image_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"status": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"status_msg": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"import_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"created": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"tag_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"comment": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"import_uuid": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourcePluginImageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

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

/*
func resourcePluginImageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
*/

func resourcePluginImageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	var diags diag.Diagnostics

	imageId := d.Id()

	imageResponse, err := c.GetImage(imageId)
	//_, err := c.GetImage(imageId)
	if err != nil {
		return diag.FromErr(err)
	}

	image := flattenPluginImageData(imageResponse)

	if err := d.Set("image", image); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(imageId)

	return diags
}

func resourcePluginImageUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourcePluginImageRead(ctx, d, m)
}

func resourcePluginImageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func flattenPluginImageData(imageResponse cn.ImageResponse) interface{} {
	image := make([]interface{}, len(imageResponse.Images), len(imageResponse.Images))

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

		image[0] = row

	}

	return image

}
