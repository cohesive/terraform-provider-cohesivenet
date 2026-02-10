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

// Plugin image V2 API and go client
func resourcePluginImageNew() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePluginImageCreateNew,
		ReadContext:   resourcePluginImageReadNew,
		UpdateContext: resourcePluginImageUpdateNew,
		DeleteContext: resourcePluginImageDeleteNew,
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
			if diff.Id() != "" && diff.HasChange("image_url") {
				return fmt.Errorf("plugin image_url cannot be changed as it has dependant plugin instances, delete the instances and recreate the plugin image")
			}
			return nil
		},
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
				Description: "Name of the plugin image",
			},
			"image_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "URL of the image to be imported",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of deployed image",
			},
			"command": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Start command for the plugin image. Either /opt/cohesive/container_startup.sh or /usr/bin/supervisord",
			},
			"documentation_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL of documentation for the plugin image",
			},
			"support_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL of support for the plugin image",
			},
			"catalog_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of image in the Cohesive Networks catalog",
			},
			"version": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Version of the plugin image",
			},
			"tags": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				//Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Key-value pairs of tags for the plugin image",
			},
			"metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				//Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Key-value pairs of metadata for the plugin image",
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
	// synchronize creating a plugin image
	vns3.ReqLock.Lock()
	defer vns3.ReqLock.Unlock()

	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	name := d.Get("name").(string)
	image_url := d.Get("image_url").(string)
	description := d.Get("description").(string)
	command := d.Get("command").(string)
	documentation_url := d.Get("documentation_url").(string)
	support_url := d.Get("support_url").(string)

	newImage := cn.NewInstallPluginRequest(name, image_url)
	newImage.SetDescription(description)
	newImage.SetCommand(command)
	newImage.SetDocumentationUrl(documentation_url)
	newImage.SetSupportUrl(support_url)

	// Get tags if they exist and extract correct data type
	if v, ok := d.GetOk("tags"); ok {
		tags := make(map[string]interface{})
		for key, value := range v.(map[string]interface{}) {
			// Convert string representations of booleans back to actual booleans
			if strVal, ok := value.(string); ok {
				switch strVal {
				case "true":
					tags[key] = true
				case "false":
					tags[key] = false
				default:
					if intVal, err := strconv.Atoi(strVal); err == nil {
						tags[key] = intVal
					} else if floatVal, err := strconv.ParseFloat(strVal, 64); err == nil {
						tags[key] = floatVal
					} else {
						// Keep as string
						tags[key] = strVal
					}
				}
			} else {
				tags[key] = value
			}
		}
		newImage.SetTags(tags)
	}

	// Get metadata extract correct data types
	if v, ok := d.GetOk("metadata"); ok {
		metadata := make(map[string]interface{})
		for key, value := range v.(map[string]interface{}) {
			// Convert string representations of booleans back to actual booleans
			if strVal, ok := value.(string); ok {
				switch strVal {
				case "true":
					metadata[key] = true
				case "false":
					metadata[key] = false
				default:
					// Check if it's a number
					if intVal, err := strconv.Atoi(strVal); err == nil {
						metadata[key] = intVal
					} else if floatVal, err := strconv.ParseFloat(strVal, 64); err == nil {
						metadata[key] = floatVal
					} else {
						// Keep as string
						metadata[key] = strVal
					}
				}
			} else {
				metadata[key] = value
			}
		}
		newImage.SetMetadata(metadata)
	}

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

func resourcePluginImageUpdateNew(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

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
	imageId := int32(iId)

	if d.HasChange("name") ||
		d.HasChange("description") ||
		d.HasChange("command") ||
		d.HasChange("documentation_url") ||
		d.HasChange("support_url") ||
		d.HasChange("tags") ||
		d.HasChange("metadata") ||
		d.HasChange("version") {

		name := d.Get("name").(string)
		description := d.Get("description").(string)
		command := d.Get("command").(string)
		documentation_url := d.Get("documentation_url").(string)
		support_url := d.Get("support_url").(string)
		version := d.Get("version").(string)

		updatedImage := cn.NewUpdatePluginRequest()
		updatedImage.SetName(name)
		updatedImage.SetDescription(description)
		updatedImage.SetCommand(command)
		updatedImage.SetDocumentationUrl(documentation_url)
		updatedImage.SetSupportUrl(support_url)
		updatedImage.SetVersion(version)

		// Get tags if they exist and extract correct data type
		if v, ok := d.GetOk("tags"); ok {
			tags := make(map[string]interface{})
			for key, value := range v.(map[string]interface{}) {
				// Convert string representations of booleans back to actual booleans
				if strVal, ok := value.(string); ok {
					switch strVal {
					case "true":
						tags[key] = true
					case "false":
						tags[key] = false
					default:
						if intVal, err := strconv.Atoi(strVal); err == nil {
							tags[key] = intVal
						} else if floatVal, err := strconv.ParseFloat(strVal, 64); err == nil {
							tags[key] = floatVal
						} else {
							// Keep as string
							tags[key] = strVal
						}
					}
				} else {
					tags[key] = value
				}
			}
			updatedImage.SetTags(tags)
		}

		apiRequest := vns3.NetworkEdgePluginsApi.PutUpdatePluginRequest(ctx, imageId)
		apiRequest = apiRequest.UpdatePluginRequest(*updatedImage)
		_, _, err := vns3.NetworkEdgePluginsApi.PutUpdatePlugin(apiRequest)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return diags

}

func resourcePluginImageDeleteNew(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

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
	imageId := int32(iId)

	apiRequest := vns3.NetworkEdgePluginsApi.DeletePluginRequest(ctx, imageId)
	_, _, err := vns3.NetworkEdgePluginsApi.DeletePlugin(apiRequest)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
