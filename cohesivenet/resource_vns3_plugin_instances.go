package cohesivenet

import (
	"context"
	"strconv"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVns3PluginInstances() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePluginInstanceCreate,
		ReadContext:   resourcePluginInstanceRead,
		UpdateContext: resourcePluginInstanceUpdate,
		DeleteContext: resourcePluginInstanceDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"plugin_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"command": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourcePluginInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	plugin_id := d.Get("plugin_id").(string)
	description := d.Get("description").(string)
	ip_address := d.Get("ip_address").(string)
	command := d.Get("command").(string)

	in := cn.CreatePluginInstance{

		//Name:        name,
		//PluginID:    plugin_id,
		//Description: description,
		//IPAddress:   ip_address,
		//Command:     command,

		Name:        name,
		ImageUUID:   plugin_id,
		Description: description,
		Ipaddress:   ip_address,
		Command:     command,
	}

	instanceResponse, err := c.CreateInstance(&in)
	if err != nil {
		return diag.FromErr(err)
	}

	id := instanceResponse.Instance.UUID
	//fmt.Println(id)
	//d.SetId(strconv.Itoa(id))
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.SetId(id)

	//resourcePluginInstanceRead(ctx, d, m)

	return diags
}

func resourcePluginInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

/*
func resourcePluginInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	var diags diag.Diagnostics

	//instanceUuid := d.Id()

	instanceResponse, err := c.GetInstance(instanceUuid)
	if err != nil {
		return diag.FromErr(err)
	}

	instance := flattenPluginInstanceData(instanceResponse)

	d.Set("name", instance[0].(map[string]interface{})["hostname"].(string))
	d.Set("ip_address", instance[0].(map[string]interface{})["ipaddress"].(string))

	//if err := d.Set("command", instance); err != nil {
	//	return diag.FromErr(err)
	//}
	//d.SetId(routeResponse.Routes[0].ID)

	d.SetId(instanceUuid)

	return diags
}
*/
func resourcePluginInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourcePluginImageRead(ctx, d, m)
}

/*
func resourcePluginInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	d.SetId("")

	return diags
}
*/
func resourcePluginInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	var diags diag.Diagnostics

	instanceUuid := d.Id()

	err := c.DeleteInstance(instanceUuid)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

/*
func resourcePluginInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	var diags diag.Diagnostics

	image := d.Get("image").([]interface{})[0]
	imageId := image.(map[string]interface{})
	id := imageId["id"].(string)

	err := c.DeleteImage(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
*/
func flattenPluginInstanceData(instanceResponse cn.InstanceResponse) []interface{} {
	image := make([]interface{}, len(instanceResponse.Instances.Containers), len(instanceResponse.Instances.Containers))

	for _, ir := range instanceResponse.Instances.Containers {
		row := make(map[string]interface{})

		row["id"] = ir.ID
		row["image"] = ir.Image
		row["hostname"] = ir.Config.Hostname
		row["ipaddress"] = ir.NetworkSettings.Networks.CohesiveNet.IPAddress
		row["path"] = ir.Path
		//row["import_id"] = ir.ImportID
		//row["created"] = ir.Created
		//row["description"] = ir.Description
		//row["comment"] = ir.Comment
		//row["import_uuid"] = ir.ImportUUID

		image[0] = row

	}

	return image

}
