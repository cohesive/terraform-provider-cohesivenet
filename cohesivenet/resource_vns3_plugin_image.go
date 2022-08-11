package cohesivenet

import (
	"context"
	"log"
	"strconv"
	"time"

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
				ForceNew: true,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
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
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"localbuild": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"localimage": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
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
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"status": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"status_msg": &schema.Schema{
							Type:     schema.TypeInt,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"import_id": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"created": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"tag_name": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"import_uuid": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
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
		Name:        image["name"].(string),
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

	//images := flattenRouteData(imageResponse)
	/*
		highest := 0
		for _, r := range images {
			values := r.(map[string]interface{})
			value, _ := strconv.Atoi(values["id"].(string))
			if value > highest {
				highest = value
			}
		}
	*/

	//d.SetId(strconv.Itoa(highest))
	uuid := imageResponse.Import_uuid
	log.Println(uuid)

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	//d.SetId(uuid)
	//resourcePluginImageRead(ctx, d, m)

	return diags
}

func resourcePluginImageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

/*
func resourcePluginImageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	var diags diag.Diagnostics

	routeId := d.Id()

	routeResponse, err := c.GetRoute(routeId)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(routeResponse.Routes[0].ID)
	return diags
}
*/

func resourcePluginImageUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourcePluginImageRead(ctx, d, m)
}

func resourcePluginImageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	// Basically we just lie and say it was deleted.
	d.SetId("")

	return diags
}

/*
func resourcePluginImageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	var diags diag.Diagnostics

	routeId := d.Id()

	err := c.DeleteRoute(routeId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func flattenPluginImageData(routeResponse cn.RouteResponse) []interface{} {
	routes := make([]interface{}, len(routeResponse.Routes), len(routeResponse.Routes))

	i := 0
	for _, rt := range routeResponse.Routes {
		row := make(map[string]interface{})

		row["cidr"] = rt.Cidr
		row["id"] = rt.ID
		row["description"] = rt.Description
		row["advertise"] = rt.Advertise
		row["metric"] = rt.Metric
		row["enabled"] = rt.Enabled
		row["netmask"] = rt.Netmask
		row["editable"] = rt.Editable
		row["table"] = rt.Table
		row["interface"] = rt.Interface

		routes[i] = row
		i++
	}

	return routes

}
*/
