package cohesivenet

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet"
	macros "github.com/cohesive/cohesivenet-client-go/cohesivenet/macros"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
				Elem: &schema.Resource{
					Schema: getVns3AuthSchema(),
				},
			},
			"new_api_password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Sets user defined API password",
			},
			"new_ui_password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Sets user defined UI password",
			},
			"new_ui_username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sets user defined admin username",
			},
			"generate_token": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Optionally creates API token",
			},
			"token_lifetime": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     31536000, // 1 year
				Description: "Sets API token lifetime. A value > 0 will generate token",
			},
			"token_refresh": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Sets API token refresh",
			},
			"license_file": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Path to VNS3 license file",
			},
			"topology_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Sets VNS3 topolgy name",
			},
			"controller_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sets VNS3 controller name",
			},
			"license_params": &schema.Schema{
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Nested block of configurable VNS3 license parameters",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Sets VNS3 overlay subnet",
						},
						"controllers": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies number of controllers in topology",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"asns": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Sets VNS3 default ASNs",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"controller_vip": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Sets VNS3 VIP",
						},
						"clients": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Sets VNS3 overlay client addresses",
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
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
				Description: "Sets VNS3 controllers peer id",
			},
			"topology_checksum": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Checksum",
			},
			"keyset_checksum": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Keyset checksum",
			},
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Token",
			},
			"new_auth_set": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"licensed": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"private_ip": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ip of controller",
			},
			"configuration_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Configuration id",
			},
			"instance_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance id (used for upgrade)",
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
				asns = append(asns, strconv.Itoa(asnAny.(int)))
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

func setVns3AuthIfCreated(vns3 *cn.VNS3Client, ctx context.Context, d *schema.ResourceData, m interface{}) *cn.VNS3Client {
	if newAPIPassword := d.Get("new_api_password").(string); newAPIPassword != "" {
		setVns3ClientPassword(vns3, newAPIPassword)
		// change v1 client
		client_v1, _ := getV1Client(ctx, d, m)
		client_v1.Password = newAPIPassword
		vns3.Log.Debug("Updated auth password")
	}
	return vns3
}

func createVns3Auth(ctx context.Context, d *schema.ResourceData, m interface{}, vns3 *cn.VNS3Client) (string, error) {
	return setVns3AuthProps(ctx, d, m, vns3, "create")
}

func updateVns3Auth(ctx context.Context, d *schema.ResourceData, m interface{}, vns3 *cn.VNS3Client) (string, error) {
	return setVns3AuthProps(ctx, d, m, vns3, "update")
}

func setVns3AuthProps(ctx context.Context, d *schema.ResourceData, m interface{}, vns3 *cn.VNS3Client, action string) (string, error) {
	newUIPassword := d.Get("new_ui_password").(string)
	newUIUsername := d.Get("new_ui_username").(string)
	newAPIPassword := d.Get("new_api_password").(string)
	createToken := d.Get("generate_token").(bool)
	tokenLifetime := d.Get("token_lifetime").(int)
	tokenRefresh := d.Get("token_refresh").(bool)

	adminUiRequest := cn.NewUpdateAdminUISettingsRequest()
	shouldUpdateUi := false
	vns3.Log.Debug(fmt.Sprintf("should change ui pass: %+v", d.HasChange("new_ui_password")))
	if (d.HasChange("new_ui_password") && newUIPassword != "") || (action == "create" && newUIPassword != "") {
		adminUiRequest.SetAdminPassword(newUIPassword)
		shouldUpdateUi = true
	}

	vns3.Log.Debug(fmt.Sprintf("should change ui username: %+v", d.HasChange("new_ui_username")))
	if d.HasChange("new_ui_username") && newUIUsername != "" || (action == "create" && newUIUsername != "") {
		adminUiRequest.SetAdminUsername(newUIUsername)
		shouldUpdateUi = true
	}

	vns3.Log.Debug(fmt.Sprintf("should update ui : %+v", shouldUpdateUi))
	if shouldUpdateUi {
		vns3.Log.Debug("Setting new UI authentication")
		uiApiRequest := vns3.ConfigurationApi.PutUpdateAdminUiRequest(ctx)
		uiApiRequest = uiApiRequest.UpdateAdminUISettingsRequest(*adminUiRequest)
		_, _, err := vns3.ConfigurationApi.PutUpdateAdminUi(uiApiRequest)
		if err != nil {
			return "UI", err
		}
	}

	vns3.Log.Debug(fmt.Sprintf("should update new_api_password : %+v", d.HasChange("new_api_password")))
	if d.HasChange("new_api_password") && newAPIPassword != "" || (action == "create" && newAPIPassword != "") {
		vns3.Log.Debug("Setting new API authentication")
		updatePassword := cn.NewUpdatePasswordRequest()
		updatePassword.SetPassword(newAPIPassword)
		apiPsSetRequest := vns3.ConfigurationApi.PutUpdateApiPasswordRequest(ctx)
		apiPsSetRequest = apiPsSetRequest.UpdatePasswordRequest(*updatePassword)
		vns3.Log.Debug("before Setting new API authentication")
		_, _, err := vns3.ConfigurationApi.PutUpdateApiPassword(apiPsSetRequest)
		vns3.Log.Debug("after Setting new API authentication")
		if err != nil {
			vns3.Log.Debug(fmt.Sprintf("err Setting new API authentication", err))
			return "API", err
		}
		setVns3AuthIfCreated(vns3, ctx, d, m)

		// if the current client is using basic auth with root api password then we need
		// to update the client
		currentClientConfig := vns3.GetConfig()
		currentAuthType := *(currentClientConfig.AuthType)
		if currentAuthType == cn.ContextBasicAuth {
			setVns3AuthIfCreated(vns3, ctx, d, m)
		}
	}

	vns3.Log.Debug(fmt.Sprintf("should update generate_token : %+v", d.HasChange("generate_token")))
	if d.HasChange("generate_token") && createToken || (action == "create" && createToken) {
		vns3.Log.Debug("Generating new API token")
		hasLifetime := tokenLifetime != 0
		hasRefresh := tokenRefresh
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
		d.Set("generate_token", false)
	}

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

	if !ctrlNameExists {
		controllerName = "ctrl"
	}

	// wait for a while if still coming up
	_, err := vns3.ConfigurationApi.WaitForApi(&ctx, 60*20, 3, 5)
	if err != nil {
		host, _ := vns3.GetConfig().ServerURL(0, map[string]string{})
		return diag.FromErr(fmt.Errorf("VNS3 is not available [host=%v] %v", host, err))
	}

	// Begin configuration
	vns3.Log.Debug(fmt.Sprintf("keysetparams config %+v", keysetParams))
	keysetSource, hasSource := keysetParams["source"]
	hasSource = hasSource && keysetSource != ""
	_licenseFile, hasLicense := d.GetOk("license_file")
	// if license is provided, we are configuring new VNS3 topology. license required
	if hasLicense {
		licenseFile := _licenseFile.(string)
		licenseParamsRequest, licenseParamsError := buildLicenseParamsRequest(d)
		if licenseParamsError != nil {
			return diag.FromErr(licenseParamsError)
		}

		keysetParamsRequest := cn.SetKeysetParamsRequest{
			Token: keysetToken,
		}

		setupReq := macros.SetupRequest{
			TopologyName:   topologyName,
			ControllerName: controllerName,
			LicenseParams:  licenseParamsRequest,
			LicenseFile:    licenseFile,
			PeerId:         peerId,
			KeysetParams:   keysetParamsRequest,
			WaitTimeout:    60 * 20,
			KeysetTimeout:  60 * 20,
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
		if !hasSource {
			return diag.FromErr(fmt.Errorf("license_file or keyset.source is required to configure VNS3"))
		}
		source := keysetSource.(string)
		_, err := macros.FetchKeysetFromSource(vns3, source, keysetToken, 60*20)
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

	// now lets set the password
	step, authErr := createVns3Auth(ctx, d, m, vns3)
	if authErr != nil {
		errMessage := fmt.Sprintf("Error updating %v authentication: %v", step, authErr.Error())
		vns3.Log.Error(errMessage)
		return diag.FromErr(fmt.Errorf(errMessage))
	}

	configDetail, _, _ := vns3.ConfigurationApi.GetConfig(vns3.ConfigurationApi.GetConfigRequest(ctx))
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

	d.SetId(strconv.Itoa(int(configData.GetManagerId())))

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

	//vns3 = setVns3AuthIfCreated(vns3, d)
	configDetail, _, err := vns3.ConfigurationApi.GetConfig(vns3.ConfigurationApi.GetConfigRequest(ctx))
	if err != nil {
		return diag.FromErr(fmt.Errorf("VNS3 Config check error: %+v", err))
	}
	configData := configDetail.GetResponse()
	topologyChecksum := configData.GetTopologyChecksum()
	d.Set("topology_checksum", topologyChecksum)
	d.Set("licensed", configData.GetLicensed())
	d.Set("private_ip", configData.GetPrivateIpaddress())

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

	var diags diag.Diagnostics
	vns3, clienterror := getVns3Client(ctx, d, m)
	if clienterror != nil {
		return diag.FromErr(clienterror)
	}

	step, authErr := updateVns3Auth(ctx, d, m, vns3)

	if authErr != nil {
		errMessage := fmt.Sprintf("Error updating %v authentication: %v", step, authErr.Error())
		vns3.Log.Error(errMessage)
		return diag.FromErr(fmt.Errorf(errMessage))
	}

	d.Set("generate_token", false) //for backwards compatibility

	if !d.HasChange("configuration_id") {
		vns3.Log.Info(fmt.Sprintf("specified VNS3 config resource cannot be updated yet coming soon %v", d.Get("controller_name").(string)))
		//notsupportederror := fmt.Errorf("specified VNS3 config resource cannot be updated yet, coming soon")
		//return diag.FromErr(notsupportederror)
		return diags
	}

	// cache the password for later after new controller i sup
	config := vns3.GetConfig()
	password := (*config.Auth).(cn.BasicAuth).Password

	// set the instance id as password for new controller instance
	setVns3ClientPassword(vns3, d.Get("instance_id").(string))

	// wait for new controller instance to come up
	_, wait_err := vns3.ConfigurationApi.WaitForApi(&ctx, 60*20, 3, 5)
	if wait_err != nil {
		host, _ := vns3.GetConfig().ServerURL(0, map[string]string{})
		return diag.FromErr(fmt.Errorf("VNS3 is not available [host=%v] %v", host, wait_err))
	}

	// take snashop of current controller
	source := d.Get("private_ip").(string)
	vns3.Log.Info(fmt.Sprintf("Init controller from source: %+v", source))

	// password is used to connect to "other" controller
	_, err := macros.InitControllerFromSource(vns3, source, password, 60*20)
	if err != nil {
		return diag.FromErr(fmt.Errorf("VNS3 Setup error [fetch snapshot] %+v", err))
	}

	// set the password back (from instance id) into vns3 client
	setVns3ClientPassword(vns3, password)
	resourceConfigRead(ctx, d, m)
	return diags
}

func resourceConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	// Basically we just lie and say it was deleted.
	d.SetId("")
	tflog.Warn(ctx, "VNS3 Configuration cannot be deleted. To re-configure VNS3 you must either 1) delete the instance and configure a VNS3 or 2) manually reset VNS3 by going to https://vns3-host:8000/reset_defaults")

	return diags
}
