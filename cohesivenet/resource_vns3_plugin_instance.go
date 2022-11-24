package cohesivenet

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//Plugin instance V2 API and go client
func resourceVns3PluginInstanceNew() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePluginInstanceCreateNew,
		ReadContext:   resourcePluginInstanceReadNew,
		UpdateContext: resourcePluginInstanceUpdateNew,
		DeleteContext: resourcePluginInstanceDeleteNew,
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
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of instance",
			},
			"plugin_id": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of plugin image",
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
			"plugin_config": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Plugin instance configuration file",
			},
		},
	}
}

func parseInstanceResponseId(pluginInstance cn.PluginInstance) string {
	iId := pluginInstance.GetId()
	instanceId := int(iId)
	instanceIdString := strconv.Itoa(instanceId)

	return instanceIdString
}

func resourcePluginInstanceCreateNew(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}
	// synchronize creating a plugin image
	vns3.ReqLock.Lock()
	defer vns3.ReqLock.Unlock()

	name := d.Get("name").(string)
	plugin_id := int32(d.Get("plugin_id").(int))
	description := d.Get("description").(string)
	ip_address := d.Get("ip_address").(string)
	command := d.Get("command").(string)

	newInstance := cn.NewStartPluginInstanceRequest(name, plugin_id)
	newInstance.SetDescription(description)
	newInstance.SetIpAddress(ip_address)
	newInstance.SetCommand(command)

	apiRequest := vns3.NetworkEdgePluginsApi.StartPluginInstanceRequest(ctx)
	apiRequest = apiRequest.StartPluginInstanceRequest(*newInstance)
	detail, _, err := vns3.NetworkEdgePluginsApi.StartPluginInstance(apiRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	instanceData := detail.GetResponse()
	instanceIdString := parseInstanceResponseId(instanceData)
	d.SetId(instanceIdString)

	Id := d.Id()
	iId, _ := strconv.Atoi(Id)
	instanceId := int32(iId)
	instanceConfig := vns3.NetworkEdgePluginsApi.GetPluginInstanceConfigContentRequest(ctx, instanceId, "0")
	instanceDetail, _, err := vns3.NetworkEdgePluginsApi.GetPluginInstanceConfigContent(instanceConfig)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println(instanceDetail)

	///all update to set new config
	resourcePluginInstanceUpdateNew(ctx, d, m)
	/*
		//Set new config
		//pluginConfig := d.Get("plugin_config").(string)
		newInstanceConfig := cn.NewUpdateFileContentRequest(pluginConfig)
		request := vns3.NetworkEdgePluginsApi.UpdatePluginInstanceConfigContentRequest(ctx, instanceId, "0")
		request = request.UpdateFileContentRequest(*newInstanceConfig)
		configDetail, _, err := vns3.NetworkEdgePluginsApi.UpdatePluginInstanceConfigFileContent(request)
		if err != nil {
			return diag.FromErr(err)
		}

		config := configDetail.GetResponse()
		d.Set("plugin_config", config)
	*/
	resourcePluginInstanceReadNew(ctx, d, m)

	return diags
}

func resourcePluginInstanceReadNew(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	Id := d.Id()
	iId, _ := strconv.Atoi(Id)
	instanceId := int32(iId)
	vns3.Log.Info(fmt.Sprintf("Reading Instance Id %v", string(instanceId)))
	getInstanceRequest := vns3.NetworkEdgePluginsApi.GetPluginInstanceRequest(ctx, instanceId)
	detail, httpResponse, err := vns3.NetworkEdgePluginsApi.GetPluginInstance(getInstanceRequest)
	if err != nil {
		if httpResponse.StatusCode == 404 {
			d.SetId("")
			return diags
		} else {
			return diag.FromErr(fmt.Errorf("VNS3 GET Plugin Instance error: %+v", err))
		}
	}

	instance := detail.GetResponse()
	instanceIdString := parseInstanceResponseId(instance)
	d.SetId(instanceIdString)

	return diags
}

func resourcePluginInstanceUpdateNew(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	Id := d.Id()
	iId, _ := strconv.Atoi(Id)
	instanceId := int32(iId)

	if d.HasChange("plugin_config") {
		pluginConfig := d.Get("plugin_config").(string)
		newInstanceConfig := cn.NewUpdateFileContentRequest(pluginConfig)
		request := vns3.NetworkEdgePluginsApi.UpdatePluginInstanceConfigContentRequest(ctx, instanceId, "0")
		request = request.UpdateFileContentRequest(*newInstanceConfig)
		configDetail, _, err := vns3.NetworkEdgePluginsApi.UpdatePluginInstanceConfigFileContent(request)
		if err != nil {
			return diag.FromErr(err)
		}

		config := configDetail.GetResponse()
		d.Set("plugin_config", config)
		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	resourcePluginInstanceReadNew(ctx, d, m)
	return diags

}

func resourcePluginInstanceDeleteNew(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}
	// synchronize creating a plugin image
	vns3.ReqLock.Lock()
	defer vns3.ReqLock.Unlock()

	Id := d.Id()
	iId, _ := strconv.Atoi(Id)
	instanceId := int32(iId)

	deleteInstance := cn.NewDeletePluginInstanceRequest()
	deleteInstance.SetForce(true)

	apiRequest := vns3.NetworkEdgePluginsApi.DeletePluginInstanceRequest(ctx, instanceId)
	apiRequest = apiRequest.DeletePluginInstanceRequest(*deleteInstance)
	_, _, err := vns3.NetworkEdgePluginsApi.DeletePluginInstance(apiRequest)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
