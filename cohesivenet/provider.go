package cohesivenet

import (
	"fmt"
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

	if token != nil && token != "" {
		auth := context.WithValue(context.Background(), cohesivenet.ContextAccessToken, token)
	} else {
		auth := context.WithValue(context.Background(), cohesivenet.ContextBasicAuth, cohesivenet.BasicAuth{
			UserName: username,
			Password: password,
		})
	}

    vns3 := cohesivenet.NewVNS3Client(cohesivenet.NewConfiguration(host), cohesivenet.ClientParams{
        Timeout: 10,
        TLS: false,
    })

	c := make(map[string]interface{})

	c["vns3_auth"] = auth
	c["vns3"] = vns3

	return c, diags

	// // cohesivenet.VNS3Client

    // req := vns3.ConfigurationApi.GetConfig(auth)

	// if (username != "") && (password != "") {
	// 	c, err := cn.NewClient(&username, &password, &token, &hostUrl)
	// 	if err != nil {
	// 		return nil, diag.FromErr(err)
	// 	}
	// 	return c, diags
	// }
	// c, err := cn.NewClient(nil, nil, nil, nil)
	// if err != nil {
	// 	return nil, diag.FromErr(err)
	// }

	// return c, diags
}
