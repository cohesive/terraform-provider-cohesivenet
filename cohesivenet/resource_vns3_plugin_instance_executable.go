package cohesivenet

import (
	"context"
	"fmt"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Plugin instance V2 API and go client
func resourceVns3PluginInstanceExecutable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePluginInstanceExecutableCreate,
		ReadContext:   resourcePluginInstanceExecutableRead,
		UpdateContext: resourcePluginInstanceExecutableUpdate,
		DeleteContext: resourcePluginInstanceExecutableDelete,
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
			"instance_id": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "ID of plugin instance the executable is running in",
			},
			"command": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Name of command to run in the plugin executable",
			},
			"executable_path": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Set path to executable inside plugin instance",
			},
			"timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Set timeout for command execution",
			},
			"output": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output from command execution",
			},
			"error": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Error response from command execution",
			},
			"timed_out": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Timeout response from command execution",
			},
			"failed": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Failed response from command execution",
			},
			"run_count": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Increment this number to re-run the command",
			},
			"last_executed": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp of last command execution",
			},
		},
	}
}

func resourcePluginInstanceExecutableCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}
	// synchronize creating a plugin image
	vns3.ReqLock.Lock()
	defer vns3.ReqLock.Unlock()

	command := d.Get("command").(string)
	executable_path := d.Get("executable_path").(string)
	instance_id := int32(d.Get("instance_id").(int))
	timeout := int32(d.Get("timeout").(int))

	newInstanceExe := cn.NewRunPluginInstanceCommandRequest(command)
	newInstanceExe.SetExecutablePath(executable_path)
	newInstanceExe.SetTimeout(timeout)

	apiRequest := vns3.NetworkEdgePluginsApi.RunPluginInstanceExecutableCommandRequest(ctx, instance_id)
	apiRequest = apiRequest.RunPluginInstanceCommandRequest(*newInstanceExe)
	apiResponse, _, err := vns3.NetworkEdgePluginsApi.RunPluginInstanceExecutableCommand(apiRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	responseData := apiResponse.GetResponse()

	if output := responseData.GetOutput(); output != "" {
		d.Set("output", output)
	}

	if timeout, ok := responseData.GetTimeoutOk(); ok {
		d.Set("timed_out", *timeout)
	}

	if failed, ok := responseData.GetFailedOk(); ok {
		d.Set("failed", *failed)
	}

	if errorMsg := responseData.GetError(); errorMsg != "" {
		d.Set("error", errorMsg)
	}

	// Only return error if the command actually failed
	if failed, ok := responseData.GetFailedOk(); ok && *failed {
		if errorMsg := responseData.GetError(); errorMsg != "" {
			return diag.FromErr(fmt.Errorf("plugin command failed: %s", errorMsg))
		}
		return diag.FromErr(fmt.Errorf("plugin command failed"))
	}

	d.Set("last_executed", time.Now().Format(time.RFC3339))
	d.SetId(time.Now().Format(time.RFC3339))

	return diags
}

func resourcePluginInstanceExecutableRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	//No-op nothing to read....

	return diags
}

func resourcePluginInstanceExecutableUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}
	// synchronize creating a plugin image
	vns3.ReqLock.Lock()
	defer vns3.ReqLock.Unlock()

	if d.HasChange("run_count") {
		command := d.Get("command").(string)
		executable_path := d.Get("executable_path").(string)
		instance_id := int32(d.Get("instance_id").(int))
		timeout := int32(d.Get("timeout").(int))

		newInstanceExe := cn.NewRunPluginInstanceCommandRequest(command)
		newInstanceExe.SetExecutablePath(executable_path)
		newInstanceExe.SetTimeout(timeout)

		apiRequest := vns3.NetworkEdgePluginsApi.RunPluginInstanceExecutableCommandRequest(ctx, instance_id)
		apiRequest = apiRequest.RunPluginInstanceCommandRequest(*newInstanceExe)
		apiResponse, _, err := vns3.NetworkEdgePluginsApi.RunPluginInstanceExecutableCommand(apiRequest)
		if err != nil {
			return diag.FromErr(err)
		}

		responseData := apiResponse.GetResponse()

		if output := responseData.GetOutput(); output != "" {
			d.Set("output", output)
		}

		if timeout, ok := responseData.GetTimeoutOk(); ok {
			d.Set("timed_out", *timeout)
		}

		if failed, ok := responseData.GetFailedOk(); ok {
			d.Set("failed", *failed)
		}

		if errorMsg := responseData.GetError(); errorMsg != "" {
			d.Set("error", errorMsg)
		}

		// Only return error if the command actually failed
		if failed, ok := responseData.GetFailedOk(); ok && *failed {
			if errorMsg := responseData.GetError(); errorMsg != "" {
				return diag.FromErr(fmt.Errorf("plugin command failed: %s", errorMsg))
			}
			return diag.FromErr(fmt.Errorf("plugin command failed"))
		}

		d.Set("last_executed", time.Now().Format(time.RFC3339))

	}

	d.SetId(time.Now().Format(time.RFC3339))

	return diags

}

func resourcePluginInstanceExecutableDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	/*
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
	*/
	d.SetId("")

	return diags
}
