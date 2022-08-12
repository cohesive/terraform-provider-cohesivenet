package cohesivenet

import (
	"context"
	"fmt"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet"
	cnv1 "github.com/cohesive/cohesivenet-client-go/cohesivenet/v1"
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
		},
		ResourcesMap: map[string]*schema.Resource{
			"cohesivenet_vns3_ipsec_endpoints":  resourceEndpoints(),
			"cohesivenet_vns3_routes":           resourceRoutes(),
			"cohesivenet_vns3_firewall_rules":   resourceRules(),
			"cohesivenet_vns3_config":           resourceVns3Config(),
			"cohesivenet_vns3_ipsec_ebpg_peers": resourceEbgp(),
			"cohesivenet_vns3_plugin_images":    resourcePluginImage(),
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

// func commonAuthSchema() *schema.Schema {
// 	return &schema.Schema{
// 		Type:     schema.TypeSet,
// 		MaxItems: 1,
// 		Optional: true,
// 		Elem:     &schema.Resource{
// 			Schema: map[string]*schema.Schema{
// 				"host": &schema.Schema{
// 					Type:    schema.TypeString,
// 					Optional: true,
// 				},
// 				"username": &schema.Schema{
// 					Type:    schema.TypeString,
// 					Optional: true,
// 				},
// 				"password": &schema.Schema{
// 					Type:    schema.TypeString,
// 					Optional: true,
// 				},
// 				"token": &schema.Schema{
// 					Type:    schema.TypeString,
// 					Optional: true,
// 				},
// 			},
// 		},
// 	}
// }

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	token := d.Get("token").(string)
	host := d.Get("host").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	fmt.Printf("%+v\n", username)
	fmt.Printf("%+v\n", password)

	var cfg *cn.Configuration
	if token != "" {
		cfg = cn.NewConfigurationWithAuth(host, cn.ContextAccessToken, token)
	} else {
		cfg = cn.NewConfigurationWithAuth(host, cn.ContextBasicAuth, cn.BasicAuth{
			UserName: username,
			Password: password,
		})
	}
	vns3 := cn.NewVNS3Client(cfg, cn.ClientParams{
		Timeout: 3,
		TLS:     false,
	})

	Logger := NewLogger(ctx)
	vns3.Log = Logger

	meta := make(map[string]interface{})

	clientv1, err := cnv1.NewClient(&username, &password, &token, &host)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	meta["clientv1"] = clientv1
	meta["vns3"] = vns3

	return meta, diags
}
