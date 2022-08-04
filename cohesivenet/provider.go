package cohesivenet

import (
	"fmt"
	"context"

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
			"cohesivenet_endpoints": resourceEndpoints(),
			"cohesivenet_routes":    resourceRoutes(),
			"cohesivenet_firewall":  resourceRules(),
			"cohesivenet_vns3_config":  resourceVns3Config(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"cohesivenet_endpoints":         dataSourceEndpoints(),
			"cohesivenet_config":            dataSourceConfig(),
			"cohesivenet_container_network": dataSourceContainerNetwork(),
			"cohesivenet_routes":            dataSourceRoutes(),
			"cohesivenet_firewall":          dataSourceFirewall(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

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
        TLS: false,
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
