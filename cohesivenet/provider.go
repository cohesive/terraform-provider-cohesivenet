package cohesivenet

import (
	"context"

	cn "github.com/cohesive/cohesivenet-client-go"
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
			"hosturl": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CN_HOSTURL", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cohesivenet_endpoints": resourceEndpoints(),
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
	hostUrl := d.Get("hosturl").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if (username != "") && (password != "") {
		c, err := cn.NewClient(&username, &password, &token, &hostUrl)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		return c, diags
	}
	c, err := cn.NewClient(nil, nil, nil, nil)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return c, diags
}
