package cohesivenet

import (
	"context"
	"os"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CN_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CN_PASSWORD", nil),
			},
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CN_TOKEN", nil),
			},
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CN_HOST", nil),
			},
			"verify_cert": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CN_VERIFY_CERT", false),
				Description: "Verify the VNS3 controller's TLS certificate. Defaults to false to preserve existing behaviour (new controllers present a self-signed certificate); the default will change to true in 2.0. Set to true once a trusted certificate is installed on the controller; while false, the API token or password can be intercepted by a man-in-the-middle.",
			},
			"vns3": &schema.Schema{
				Type:       schema.TypeSet,
				MaxItems:   1,
				Optional:   true,
				Deprecated: vns3BlockDeprecation,
				Elem: &schema.Resource{
					Schema: getVns3AuthSchema(),
				},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cohesivenet_vns3_ipsec_endpoints":            resourceEndpoints(),
			"cohesivenet_vns3_ipsec_tunnel":               resourceTunnel(),
			"cohesivenet_vns3_ipsec_traffic_pair":         resourceTrafficPair(),
			"cohesivenet_vns3_routes":                     resourceRoutes(),
			"cohesivenet_vns3_firewall_rules":             resourceRules(),
			"cohesivenet_vns3_firewall_rule":              resourceFirewallRules(),
			"cohesivenet_vns3_firewall_fwset":             resourceFwSet(),
			"cohesivenet_vns3_ipsec_ebpg_peers":           resourceEbgp(),
			"cohesivenet_vns3_plugin_images":              resourcePluginImage(),
			"cohesivenet_vns3_plugin_image":               resourcePluginImageNew(),
			"cohesivenet_vns3_config":                     resourceVns3Config(),
			"cohesivenet_vns3_peers":                      resourceVns3Peering(),
			"cohesivenet_vns3_link":                       resourceLink(),
			"cohesivenet_vns3_plugin_instances":           resourceVns3PluginInstances(),
			"cohesivenet_vns3_plugin_instance":            resourceVns3PluginInstanceNew(),
			"cohesivenet_vns3_https_certs":                resourceHttpsCerts(),
			"cohesivenet_vns3_alert":                      resourceAlert(),
			"cohesivenet_vns3_webhook":                    resourceWebhook(),
			"cohesivenet_vns3_identity_vpn":               resourceIdentityVpn(),
			"cohesivenet_vns3_identity_controller":        resourceIdentityController(),
			"cohesivenet_vns3_license_upgrade":            resourceLicenseUpgrade(),
			"cohesivenet_vns3_plugin_instance_executable": resourceVns3PluginInstanceExecutable(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"cohesivenet_vns3_ipsec_endpoints":   dataSourceEndpoints(),
			"cohesivenet_vns3_config":            dataSourceConfig(),
			"cohesivenet_vns3_container_network": dataSourceContainerNetwork(),
			"cohesivenet_vns3_route":             dataSourceRoutes(),
			"cohesivenet_vns3_firewall":          dataSourceFirewall(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

/* configure VNS3 provider */
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	Logger := NewLogger(ctx)
	meta := make(map[string]interface{})

	//check if VNS3 config is defined in "vns3"
	vns3AuthSet, hasVns3Auth := d.Get("vns3").(*schema.Set)

	if hasVns3Auth && vns3AuthSet.Len() != 0 {
		vns3Auth := vns3AuthSet.List()[0].(map[string]any)
		//create v1 client
		v1client, err := generateV1Client(vns3Auth, Logger)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		meta["clientv1"] = v1client
		//create vns3 client
		vns3Client, err := generateVNS3Client(vns3Auth, Logger)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		meta["vns3"] = vns3Client
	} else {
		vns3Auth := make(map[string]interface{})
		vns3Auth["host"] = d.Get("host").(string)
		vns3Auth["username"] = d.Get("username").(string)
		vns3Auth["password"] = d.Get("password").(string)
		vns3Auth["token"] = d.Get("token").(string)
		vns3Auth["verify_cert"] = d.Get("verify_cert").(bool)
		if vns3Auth["host"] != "" {
			// Warn when verification is off only because of the default, so users
			// make a deliberate choice. Explicitly setting verify_cert (config or env) stays silent.
			_, envSet := os.LookupEnv("CN_VERIFY_CERT")
			if shouldWarnUnverifiedDefault(d.GetRawConfig(), vns3Auth["verify_cert"].(bool), envSet) {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "TLS certificate verification is disabled by default",
					Detail: "The provider does not verify the VNS3 controller's TLS certificate unless verify_cert = true is set, " +
						"so the API token or password could be intercepted by a man-in-the-middle. Once a trusted certificate is installed " +
						"on the controller, set verify_cert = true. To acknowledge and silence this warning, set verify_cert = false. " +
						"The default will change to true in version 2.0.",
				})
			}
			//create v1 client
			v1client, err := generateV1Client(vns3Auth, Logger)
			if err != nil {
				return nil, diag.FromErr(err)
			}
			meta["clientv1"] = v1client
			//create vns3 client
			vns3Client, err := generateVNS3Client(vns3Auth, Logger)
			if err != nil {
				return nil, diag.FromErr(err)
			}
			meta["vns3"] = vns3Client
		}
	}
	return meta, diags
}

// shouldWarnUnverifiedDefault reports whether verification is off only because
// of the provider default: verify is false, the env var is not set and the
// attribute is absent (null) in the configuration.
func shouldWarnUnverifiedDefault(raw cty.Value, verify bool, envSet bool) bool {
	if verify || envSet {
		return false
	}
	if raw.IsNull() || !raw.IsKnown() || !raw.Type().IsObjectType() || !raw.Type().HasAttribute("verify_cert") {
		return false
	}
	return raw.GetAttr("verify_cert").IsNull()
}
