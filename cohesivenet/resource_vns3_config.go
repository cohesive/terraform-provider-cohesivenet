package cohesivenet

import (
	"context"
	"fmt"
	"strconv"
	"time"

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
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"apitoken": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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

func getVns3Client(ctx context.Context, d *schema.ResourceData, m interface{}) (*cn.VNS3Client, error) {
	host, _ := d.Get("host").(string)
	password, _ := d.Get("password").(string)
	token, _ := d.Get("token").(string)

	var vns3 *cn.VNS3Client
	// this is a lot of code to just determine vns3 client
	if host != "" || password != "" || token != "" {
		invalid := host == "" || (password == "" && token == "")
		if invalid {
			return nil, fmt.Errorf("host and auth is required if host, password or token are passed")
		}

		var cfg *cn.Configuration
		if token != "" {
			cfg = cn.NewConfigurationWithAuth(host, cn.ContextAccessToken, token)
		} else {
			cfg = cn.NewConfigurationWithAuth(host, cn.ContextBasicAuth, cn.BasicAuth{
				UserName: "api",
				Password: password,
			})
		}

		vns3 = cn.NewVNS3Client(cfg, cn.ClientParams{
			Timeout: 3,
			TLS:     false,
		})
		Logger := NewLogger(ctx)
		vns3.Log = Logger

	} else {
		vns3_ := m.(map[string]interface{})["vns3"].(cn.VNS3Client)
		vns3 = &vns3_
	}

	return vns3, nil
}

func resourceConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	licenseParamsSet, hasParams := d.Get("license_params").(*schema.Set)

	if hasParams {
		licenseParams := licenseParamsSet.List()[0]
		vns3.Log.Info(fmt.Sprintf("License params passed %+v", licenseParams))
	}

	// Keyset params are required so will eist.
	keysetParams := d.Get("keyset_params").(*schema.Set).List()[0].(map[string]interface{})

	licenseParamsRequest := cn.NewSetLicenseParametersRequest(true)
	// TODO, set other keyset params if there
	keysetParamsRequest := cn.SetKeysetParamsRequest{
		Token: keysetParams["token"].(string),
	}

	vns3.Log.Debug(fmt.Sprintf("keysetparams config %+v", keysetParams))

	topologyName := d.Get("topology_name").(string)
	controllerName, ctrlNameExists := d.Get("controller_name").(string)
	if !ctrlNameExists {
		controllerName = "ctrl"
	}

	setupReq := macros.SetupRequest{
		TopologyName:   topologyName,
		ControllerName: controllerName,
		LicenseParams:  licenseParamsRequest,
		//LicenseFile: "/Users/benplatta/code/cohesive/vns3-functional-testing/test-assets/license.txt",
		LicenseFile:   "/Users/scott/vfuc-test-sme-license.txt",
		PeerId:        1,
		KeysetParams:  keysetParamsRequest,
		WaitTimeout:   60 * 5,
		KeysetTimeout: 60 * 5,
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
