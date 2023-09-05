package cohesivenet

import (
	"context"

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
			"vns3": &schema.Schema{
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: getVns3AuthSchema(),
				},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cohesivenet_vns3_ipsec_endpoints":    resourceEndpoints(),
			"cohesivenet_vns3_ipsec_tunnel":       resourceTunnel(),
			"cohesivenet_vns3_ipsec_traffic_pair": resourceTrafficPair(),
			"cohesivenet_vns3_routes":             resourceRoutes(),
			"cohesivenet_vns3_firewall_rules":     resourceRules(),
			"cohesivenet_vns3_firewall_rule":      resourceFirewallRules(),
			"cohesivenet_vns3_firewall_fwset":     resourceFwSet(),
			"cohesivenet_vns3_ipsec_ebpg_peers":   resourceEbgp(),
			"cohesivenet_vns3_plugin_images":      resourcePluginImage(),
			"cohesivenet_vns3_plugin_image":       resourcePluginImageNew(),
			"cohesivenet_vns3_config":             resourceVns3Config(),
			"cohesivenet_vns3_peers":              resourceVns3Peering(),
			"cohesivenet_vns3_link":               resourceLink(),
			"cohesivenet_vns3_plugin_instances":   resourceVns3PluginInstances(),
			"cohesivenet_vns3_plugin_instance":    resourceVns3PluginInstanceNew(),
			"cohesivenet_vns3_https_certs":        resourceHttpsCerts(),
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
		if vns3Auth["host"] != "" {
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
