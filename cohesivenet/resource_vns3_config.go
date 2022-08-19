package cohesivenet

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
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
				Elem:     &schema.Resource{
					Schema: getVns3AuthSchema(),
				},
			},
			"new_api_password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Sensitive: true,
			},
			"new_ui_password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Sensitive: true,
			},
			"new_ui_username": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"generate_token": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"token_lifetime": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"token_refresh": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"license_file": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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
			"token": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"new_auth_set": &schema.Schema{
				Type:     schema.TypeBool,
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

		subnetL := licenseParams["subnet"]
		hasSubnet := subnetL != ""
		if hasSubnet {
			licenseParamsRequest.SetSubnet(subnetL.(string))
			hasCustomParams = true
		} else {
			missingParams = append(missingParams, "subnet")
		}

		controllersL := licenseParams["controllers"].([]any)
		hasControllers := len(controllersL) > 0
		if hasControllers {
			controllers := []string{}
			for _, ipAny := range controllersL {
				controllers = append(controllers, ipAny.(string))
			}

			licenseParamsRequest.SetManagers(
				strings.Join(controllers, " "),
			)
			hasCustomParams = true
		} else {
			missingParams = append(missingParams, "controllers")
		}

		asnsL := licenseParams["asns"].([]any)
		hasAsns := len(asnsL) > 0
		if hasAsns {
			asns := []string{}
			for _, asnAny := range asnsL {
				asns = append(asns, asnAny.(string))
			}
			licenseParamsRequest.SetAsns(
				strings.Join(asns, " "),
			)
			hasCustomParams = true
		} else {
			missingParams = append(missingParams, "asns")
		}

		controllerVip := licenseParams["controller_vip"].(string)
		hasVip := controllerVip != ""
		if hasVip {
			licenseParamsRequest.SetMyManagerVip(controllerVip)
			hasCustomParams = true
		} else {
			missingParams = append(missingParams, "controller_vip")
		}

		clientsL := licenseParams["clients"].([]any)
		hasClients := len(clientsL) > 0
		
		if hasClients {
			clients := []string{}
			for _, clientIpAny := range clientsL {
				clients = append(clients, clientIpAny.(string))
			}
			licenseParamsRequest.SetClients(strings.Join(clients, ","))
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

func setVns3AuthIfCreated(vns3 *cn.VNS3Client, d *schema.ResourceData) *cn.VNS3Client {
	_id := d.Id()
	if _id != "" {
		if newAPIPassword := d.Get("new_api_password").(string); newAPIPassword != "" {
			setVns3ClientPassword(vns3, newAPIPassword)
		}
	}
	return vns3
}

func updateVns3Auth(ctx context.Context, d *schema.ResourceData, vns3 *cn.VNS3Client) (string, error) {
	newUIPassword := d.Get("new_ui_password").(string)
	newUIUsername := d.Get("new_ui_username").(string)
	newAPIPassword := d.Get("new_api_password").(string)
	createToken := d.Get("generate_token").(bool)
	tokenLifetime := d.Get("token_lifetime").(int)
	tokenRefresh := d.Get("token_refresh").(bool)

	adminUiRequest := cn.NewUpdateAdminUISettingsRequest()
	shouldUpdateUi := false
	if newUIPassword != "" {
		adminUiRequest.SetAdminPassword(newUIPassword)
		shouldUpdateUi = true
	}

	if newUIUsername != "" {
		adminUiRequest.SetAdminUsername(newUIUsername)
		shouldUpdateUi = true
	}

	if shouldUpdateUi {
		vns3.Log.Debug("Setting new UI authentication")
		uiApiRequest := vns3.ConfigurationApi.PutUpdateAdminUiRequest(ctx)
		uiApiRequest = uiApiRequest.UpdateAdminUISettingsRequest(*adminUiRequest)
		_, _, err := vns3.ConfigurationApi.PutUpdateAdminUi(uiApiRequest)
		if err != nil {
			return "UI", err
		}
	}

	if newAPIPassword != "" {
		vns3.Log.Debug("Setting new API authentication")
		updatePassword := cn.NewUpdatePasswordRequest()
		updatePassword.SetPassword(newAPIPassword)
		apiPsSetRequest := vns3.ConfigurationApi.PutUpdateApiPasswordRequest(ctx)
		apiPsSetRequest = apiPsSetRequest.UpdatePasswordRequest(*updatePassword)
		_, _, err := vns3.ConfigurationApi.PutUpdateApiPassword(apiPsSetRequest)
		if err != nil {
			return "API", err
		}

		// if the current client is using basic auth with root api password then we need
		// to update the client
		currentClientConfig := vns3.GetConfig()
		currentAuthType := *(currentClientConfig.AuthType)
		if currentAuthType == cn.ContextBasicAuth {
			setVns3ClientPassword(vns3, newAPIPassword)
		}
	}

	if createToken {
		vns3.Log.Debug("Generating new API token")
		hasLifetime := tokenLifetime != 0
		hasRefresh := tokenRefresh != false
		var tokenRequest *cn.CreateAPITokenRequest
		if hasLifetime || hasRefresh {
			tokenRequest = cn.NewCreateAPITokenRequest()
			if hasLifetime {
				tokenRequest.SetExpires(int32(tokenLifetime))
			}
			if hasRefresh {
				tokenRequest.SetRefreshes(tokenRefresh)
			}
		}

		apiRequest := vns3.AccessApi.CreateApiTokenRequest(ctx)
		if tokenRequest != nil {
			apiRequest = apiRequest.CreateAPITokenRequest(*tokenRequest)
		}

		tokenResponse, _, err := vns3.AccessApi.CreateApiToken(apiRequest)
		if err != nil {
			return "token", err
		}
		accessTokenData := tokenResponse.GetResponse()

		d.Set("token", accessTokenData.Token)
	}

	// set a flag for handling auth setting
	vns3.Log.Debug("New VNS3 auth setting run. Setting flag")
	d.Set("new_auth_set", true)
	return "", nil
}

func resourceConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	// required params
	topologyName := d.Get("topology_name").(string)
	keysetParams := d.Get("keyset_params").(*schema.Set).List()[0].(map[string]interface{})
	peerId := int32(d.Get("peer_id").(int))
	keysetToken := keysetParams["token"].(string)
	controllerName, ctrlNameExists := d.Get("controller_name").(string)
	hasAuthBeenSet := d.Get("new_auth_set").(bool)
	newAPIPassword := d.Get("new_api_password").(string)
	if hasAuthBeenSet && newAPIPassword != "" {
		setVns3ClientPassword(vns3, newAPIPassword)
	}

	if !ctrlNameExists {
		controllerName = "ctrl"
	}

	// wait for a while if still coming up
	_, err := vns3.ConfigurationApi.WaitForApi(&ctx, 60*8, 3, 5)
	if err != nil {
		host, _ := vns3.GetConfig().ServerURL(0, map[string]string{})
		return diag.FromErr(fmt.Errorf("VNS3 is not available [host=%v]", host))
	}

	if !hasAuthBeenSet {
		step, authErr := updateVns3Auth(ctx, d, vns3)
		if authErr != nil {
			errMessage := fmt.Sprintf("Error updating %v authentication: %v", step, authErr.Error())
			vns3.Log.Error(errMessage)
			return diag.FromErr(fmt.Errorf(errMessage))
		}
	}

	// Begin configuration
	vns3.Log.Debug(fmt.Sprintf("keysetparams config %+v", keysetParams))

	keysetSource, hasSource := keysetParams["source"]
	hasSource = hasSource && keysetSource != ""
	// if no source provided, we are configuring new VNS3 topology. license required
	if !hasSource {
		_licenseFile, hasLicense := d.GetOk("license_file")
		if !hasLicense {
			return diag.FromErr(fmt.Errorf("license_file or keyset.source is required to configure VNS3"))
		}

		licenseFile := _licenseFile.(string)
		licenseParamsRequest, licenseParamsError := buildLicenseParamsRequest(d)
		if licenseParamsError != nil {
			return diag.FromErr(licenseParamsError)
		}

		keysetParamsRequest := cn.SetKeysetParamsRequest{
			Token: keysetToken,
		}

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

		configDetail, setupErr := macros.SetupController(vns3, setupReq)

		if setupErr != nil {
			vns3.Log.Error(fmt.Sprintf("VNS3 Setup error: %+v", setupErr))
			return diag.FromErr(fmt.Errorf("VNS3 Setup error: %+v", setupErr))
		} else {
			c := *configDetail
			d, _ := c.MarshalJSON()
			vns3.Log.Info(fmt.Sprintf("VNS3 Setup success %+v", string(d)))
		}

	} else {
		source := keysetSource.(string)
		_, err := macros.FetchKeysetFromSource(vns3, source, keysetToken, 60*10)
		if err != nil {
			return diag.FromErr(fmt.Errorf("VNS3 Setup error [fetch keyset] %+v", err))
		}
		// set topology and ctrl name
		configReq := cn.NewUpdateConfigRequest()
		configReq.SetTopologyName(topologyName)
		configReq.SetControllerName(controllerName)
		apiReq := vns3.ConfigurationApi.PutConfigRequest(ctx).UpdateConfigRequest(*configReq)
		vns3.ConfigurationApi.PutConfig(apiReq)
		// set peering
		_, peeringErr := macros.TrySetPeering(vns3, peerId)
		if peeringErr != nil {
			return diag.FromErr(fmt.Errorf("VNS3 Setup error [peering] %+v", err))
		}
	}

	configDetail, _, err := vns3.ConfigurationApi.GetConfig(vns3.ConfigurationApi.GetConfigRequest(ctx))
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

	vns3 = setVns3AuthIfCreated(vns3, d)
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
	tflog.Warn(ctx, "VNS3 Configuration cannot be deleted. To re-configure VNS3 you must either 1) delete the instance and configure a VNS3 or 2) manually reset VNS3 by going to https://vns3-host:8000/reset_defaults")

	return diags
}
