package cohesivenet

import (
	"context"
	"fmt"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)


func getVns3Client(ctx context.Context, d *schema.ResourceData, m interface{}) (*cn.VNS3Client, error) {
	vns3AuthSet, hasVns3Auth := d.Get("vns3").(*schema.Set)

	var vns3 *cn.VNS3Client

	if hasVns3Auth {
		vns3Auth := vns3AuthSet.List()[0].(map[string]any)
		vns3Host, hasHost := vns3Auth["host"];
		if !hasHost {
			return nil, fmt.Errorf("vns3 block requires host param and an authentication method")
		}

		host := vns3Host.(string)
		var cfg *cn.Configuration
		if vns3Ps, hasPs := vns3Auth["password"]; hasPs {
			vns3Username, hasUsername := vns3Auth["username"]
			username := "api"
			if hasUsername {
				username = vns3Username.(string)
			}

			cfg = cn.NewConfigurationWithAuth(host, cn.ContextBasicAuth, cn.BasicAuth{
				UserName: username,
				Password: vns3Ps.(string),
			})
		} else {
			apiToken, hasToken := vns3Auth["api_token"]
			if !hasToken {
				return nil, fmt.Errorf("vns3 block requires host param and an authentication method: either password or api_token")
			}

			cfg = cn.NewConfigurationWithAuth(host, cn.ContextAccessToken, apiToken.(string))
		}

		vns3 = cn.NewVNS3Client(cfg, cn.ClientParams{
			Timeout: 3,
			TLS: false,
		})
		Logger := NewLogger(ctx)
		vns3.Log = Logger
	} else {
		vns3_ := m.(map[string]interface{})["vns3"].(cn.VNS3Client)
		vns3 = &vns3_
	}

	return vns3, nil
}


func getVns3AuthSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"host": &schema.Schema{
			Type:    schema.TypeString,
			Optional: true,
		},
		"password": &schema.Schema{
			Type:    schema.TypeString,
			Optional: true,
		},
		"api_token": &schema.Schema{
			Type:    schema.TypeString,
			Optional: true,
		},
		"username": &schema.Schema{
			Type:    schema.TypeString,
			Optional: true,
		},
	}
}