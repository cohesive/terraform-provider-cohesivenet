package cohesivenet

import (
	"context"
	"fmt"
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
		// Currently update is not supported due to complexity.
		//UpdateContext: resourcePluginInstanceUpdate,
		DeleteContext: resourcePluginInstanceDelete,
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
				Optional:    true,
				ForceNew:    true,
				Description: "Name of instance",
			},
			"plugin_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Id of instance",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Description of instance",
			},
			"ip_address": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "IP address of deployed image",
			},
			"command": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Command used to start instance",
			},
			"environment": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Environment variables used when launching instance",
			},
		},
	}
}

func resourcePluginInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	plugin_id := d.Get("plugin_id").(string)
	description := d.Get("description").(string)
	ip_address := d.Get("ip_address").(string)
	command := d.Get("command").(string)
	environment := d.Get("environment").(string)

	in := cn.CreatePluginInstance{
		Name:        name,
		ImageUUID:   plugin_id,
		Description: description,
		Ipaddress:   ip_address,
		Command:     command,
		Environment: environment,
	}

	instanceResponse, err := c.CreateInstance(&in)
	if err != nil {
		return diag.FromErr(err)
	}

	id := instanceResponse.Instance.UUID
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.SetId(id)

	resourcePluginInstanceRead(ctx, d, m)

	return diags
}

func resourcePluginInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	var diags diag.Diagnostics

	instanceUuid := d.Id()

	instanceResponse, err := c.GetInstance(instanceUuid)
	if err != nil {
		return diag.FromErr(err)
	}

	instance := flattenPluginInstanceData(instanceResponse)

	fmt.Println(instance)

	d.SetId(instanceUuid)

	return diags
}

func resourcePluginInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	var diags diag.Diagnostics

	instanceUuid := d.Id()

	err := c.DeleteInstance(instanceUuid)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func flattenPluginInstanceData(instanceResponse cn.InstanceResponse) []interface{} {
	instance := make([]interface{}, len(instanceResponse.Instances))

	for _, ir := range instanceResponse.Instances {
		row := make(map[string]interface{})

		row["Id"] = ir.ID
		row["Image"] = ir.Image
		row["Hostname"] = ir.Hostname
		row["Ipaddress"] = ir.IPAddress
		row["Path"] = ir.Path

		instance[0] = row

	}

	return instance

}
