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

	Logger := NewLogger(ctx)
	if hasVns3Auth {
		vns3Auth := vns3AuthSet.List()[0].(map[string]any)
		vns3Host := vns3Auth["host"].(string);
		hasHost := vns3Host != ""
		if !hasHost {
			return nil, fmt.Errorf("vns3 block requires host param and an authentication method")
		}

		host := vns3Host
		var cfg *cn.Configuration
		password := vns3Auth["password"].(string)
		hasPassword := password != ""
		if hasPassword {
			vns3Username := vns3Auth["username"].(string)
			var username string
			if vns3Username != "" {
				username = vns3Username
			} else {
				username = "api"
			}

			Logger.Debug("Using Basic auth for VNS3")
			cfg = cn.NewConfigurationWithAuth(host, cn.ContextBasicAuth, cn.BasicAuth{
				UserName: username,
				Password: password,
			})
		} else {
			Logger.Debug("Using API Token auth for VNS3")
			apiToken := vns3Auth["api_token"].(string)
			hasToken := apiToken != ""
			if !hasToken {
				return nil, fmt.Errorf("vns3 block requires host param and an authentication method: either password or api_token")
			}

			cfg = cn.NewConfigurationWithAuth(host, cn.ContextAccessToken, apiToken)
		}

		vns3 = cn.NewVNS3Client(cfg, cn.ClientParams{
			Timeout: 3,
			TLS: false,
		})
	} else {
		vns3_ := m.(map[string]interface{})["vns3"].(cn.VNS3Client)
		vns3 = &vns3_
	}

	vns3.Log = Logger

	return vns3, nil
}

func setVns3ClientPassword(vns3 *cn.VNS3Client, newPassword string) *cn.VNS3Client {
	vns3.SetAuth(cn.ContextBasicAuth, cn.BasicAuth{
		UserName: "api", // will change if we support different API users
		Password: newPassword,
	})
	return vns3
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
			Sensitive:   true,
		},
		"api_token": &schema.Schema{
			Type:    schema.TypeString,
			Optional: true,
			Sensitive:   true,
		},
		"username": &schema.Schema{
			Type:    schema.TypeString,
			Optional: true,
		},
	}
}