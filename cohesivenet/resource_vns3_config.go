package cohesivenet

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"strings"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet"
	macros "github.com/cohesive/cohesivenet-client-go/cohesivenet/macros"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVns3Config() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConfigCreate,
		ReadContext:   resourceConfigRead,
		UpdateContext: resourceConfigUpdate,
		DeleteContext: resourceConfigDelete,
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
				Elem:     &schema.Resource{
					Schema: getVns3AuthSchema(),
				},
			},
			"license_file": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"topology_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"controller_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"license_params": &schema.Schema{
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet": &schema.Schema{
							Type:    schema.TypeString,
							Optional: true,
						},
						"controllers": &schema.Schema{
							Type:    schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"asns": &schema.Schema{
							Type:    schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"controller_vip": &schema.Schema{
							Type:    schema.TypeString,
							Optional: true,
						},
						"clients": &schema.Schema{
							Type:    schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"default": &schema.Schema{
							Type:     schema.TypeBool,
							Default:  true,
							Optional: true,
						},
					},
				},
			},
			"keyset_params": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"token": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"source": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"peer_id": &schema.Schema{
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},
			"topology_checksum": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"keyset_checksum": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"licensed": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func buildLicenseParamsRequest(d *schema.ResourceData) (*cn.SetLicenseParametersRequest, error) {
	licenseParamsSet, hasParams := d.Get("license_params").(*schema.Set)

	licenseParamsRequest := cn.NewSetLicenseParametersRequest(true)
	if hasParams {
		licenseParams := licenseParamsSet.List()[0].(map[string]any)
		// "subnet,controllers,asns,controller_vip,clients,default"

		hasCustomParams := false
		missingParams := []string{}

		if subnetL, hasSubnet := licenseParams["subnet"]; hasSubnet {
			licenseParamsRequest.SetSubnet(subnetL.(string))
			hasCustomParams = true
		} else {
			missingParams = append(missingParams, "subnet")
		}


		if controllersL, hasControllers := licenseParams["controllers"]; hasControllers {
			licenseParamsRequest.SetManagers(
				strings.Join(controllersL.([]string), " "),
			)
			hasCustomParams = true
		} else {
			missingParams = append(missingParams, "controllers")
		}

		if asnsL, hasAsns := licenseParams["asns"]; hasAsns {
			licenseParamsRequest.SetAsns(
				strings.Join(asnsL.([]string), " "),
			)
			hasCustomParams = true
		} else {
			missingParams = append(missingParams, "asns")
		}

		if controllerVip, hasVip := licenseParams["controller_vip"]; hasVip {
			licenseParamsRequest.SetMyManagerVip(controllerVip.(string))
			hasCustomParams = true
		} else {
			missingParams = append(missingParams, "controller_vip")
		}
		
		if clientsL, hasClients := licenseParams["clients"]; hasClients {
			licenseParamsRequest.SetClients(strings.Join(clientsL.([]string), ","))
			hasCustomParams = true
		} else {
			missingParams = append(missingParams, "clients")
		}


		if hasCustomParams {
			if len(missingParams) != 0 {
				return nil, fmt.Errorf("subnet, controllers, asns, controller_vip and clients required if default is false")
			} else {
				licenseParamsRequest.SetDefault(false)
			}
		}
	}

	return licenseParamsRequest, nil
}

func resourceConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	licenseParamsRequest, licenseParamsError := buildLicenseParamsRequest(d)
	if licenseParamsError != nil {
		return diag.FromErr(licenseParamsError)
	}

	// Keyset params are required so will eist.
	keysetParams := d.Get("keyset_params").(*schema.Set).List()[0].(map[string]interface{})

	// TODO, set other keyset params if there
	keysetParamsRequest := cn.SetKeysetParamsRequest{
		Token: keysetParams["token"].(string),
	}

	keysetSource, hasSource := keysetParams["source"]
	if hasSource {
		keysetParamsRequest.SetSource(keysetSource.(string))
	}

	vns3.Log.Debug(fmt.Sprintf("keysetparams config %+v", keysetParams))

	topologyName := d.Get("topology_name").(string)
	controllerName, ctrlNameExists := d.Get("controller_name").(string)
	if !ctrlNameExists {
		controllerName = "ctrl"
	}


	peerId := d.Get("peer_id").(int32)
	licenseFile := d.Get("license_file").(string)

    setupReq := macros.SetupRequest{
        TopologyName: topologyName,
        ControllerName: controllerName,
        LicenseParams: licenseParamsRequest,
        LicenseFile: licenseFile,
        PeerId: peerId,
        KeysetParams: keysetParamsRequest,
        WaitTimeout: 60*10,
        KeysetTimeout: 60*5,
    }

	// wait for a while if still coming up
	_, err := vns3.ConfigurationApi.WaitForApi(&ctx, 60*10, 3, 5)
	configDetail, setupErr := macros.SetupController(vns3, setupReq)

	if setupErr != nil {
		vns3.Log.Error(fmt.Sprintf("VNS3 Setup error: %+v", setupErr))
		return diag.FromErr(fmt.Errorf("VNS3 Setup error: %+v", setupErr))
	} else {
		c := *configDetail
		d, _ := c.MarshalJSON()
		vns3.Log.Info(fmt.Sprintf("VNS3 Setup success %+v", string(d)))
	}

	configData := configDetail.GetResponse()
	topologyChecksum := configData.GetTopologyChecksum()
	d.Set("topology_checksum", topologyChecksum)
	d.Set("licensed", configData.GetLicensed())

	keysetDetail, _, err := vns3.ConfigurationApi.GetKeyset(vns3.ConfigurationApi.GetKeysetRequest(ctx))

	if err != nil {
		return diag.FromErr(fmt.Errorf("VNS3 Keyset check error: %+v", err))
	}

	keysetData := keysetDetail.GetResponse()
	keysetChecksum := keysetData.GetChecksum()
	d.Set("keyset_checksum", keysetChecksum)

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

/*
func resourceConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

*/
func resourceConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	configDetail, _, err := vns3.ConfigurationApi.GetConfig(vns3.ConfigurationApi.GetConfigRequest(ctx))
	if err != nil {
		return diag.FromErr(fmt.Errorf("VNS3 Config check error: %+v", err))
	}

	configData := configDetail.GetResponse()
	topologyChecksum := configData.GetTopologyChecksum()
	d.Set("topology_checksum", topologyChecksum)
	d.Set("licensed", configData.GetLicensed())

	keysetDetail, _, err := vns3.ConfigurationApi.GetKeyset(vns3.ConfigurationApi.GetKeysetRequest(ctx))

	if err != nil {
		return diag.FromErr(fmt.Errorf("VNS3 Keyset check error: %+v", err))
	}

	keysetData := keysetDetail.GetResponse()
	keysetChecksum := keysetData.GetChecksum()
	d.Set("keyset_checksum", keysetChecksum)

	return diags
}

func resourceConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// TODO: we could allow topology name and controller name to be reset and only fail
	// when license params or keyset params change
	notsupportederror := fmt.Errorf("VNS3 config resource cannot be updated. Please redeploy a new server or reset defaults and edit terraform state")
	return diag.FromErr(notsupportederror)
}

func resourceConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	// Basically we just lie and say it was deleted.
	d.SetId("")

	return diags
}
