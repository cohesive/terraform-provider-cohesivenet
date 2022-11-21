package cohesivenet

import (
	"context"
	"fmt"
	"strings"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet"
	cnv1 "github.com/cohesive/cohesivenet-client-go/cohesivenet/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getVns3Client(ctx context.Context, d *schema.ResourceData, m interface{}) (*cn.VNS3Client, error) {
	Logger := NewLogger(ctx)
	defaultVns3Client, hasDefaultVns3Client := m.(map[string]interface{})["vns3"].(*cn.VNS3Client)
	vns3AuthSet, hasVns3Auth := d.Get("vns3").(*schema.Set)
	if hasVns3Auth && vns3AuthSet.Len() != 0 {
		vns3Auth := vns3AuthSet.List()[0].(map[string]any)
		vns3Client, err := generateVNS3Client(vns3Auth, Logger)
		if err != nil {
			return nil, err
		}
		if hasDefaultVns3Client {
			vns3Client.ReqLock = defaultVns3Client.ReqLock
		}
		return vns3Client, nil
	}
	if !hasDefaultVns3Client {
		return nil, fmt.Errorf("no vns3 configured in provider or in resource")
	}
	return defaultVns3Client, nil
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
			Type:     schema.TypeString,
			Optional: true,
		},
		"password": &schema.Schema{
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"api_token": &schema.Schema{
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"username": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"timeout": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
	}
}

func getV1Client(ctx context.Context, d *schema.ResourceData, m interface{}) (*cnv1.Client, error) {
	Logger := NewLogger(ctx)
	defaultV1Client, hasDefaultV1Client := m.(map[string]interface{})["clientv1"].(*cnv1.Client)
	//first check for vns3 override
	vns3AuthSet, hasVns3Auth := d.Get("vns3").(*schema.Set)
	if hasVns3Auth && vns3AuthSet.Len() != 0 {
		vns3Auth := vns3AuthSet.List()[0].(map[string]any)
		v1client, err := generateV1Client(vns3Auth, Logger)
		if err != nil {
			return nil, err
		}
		if hasDefaultV1Client {
			v1client.ReqLock = defaultV1Client.ReqLock
		}
		return v1client, nil
	} else {
		if !hasDefaultV1Client {
			return nil, fmt.Errorf("no vns3 configured in provider or in resource")
		}
		return defaultV1Client, nil
	}
}

//will be deprecated when all resources using Go client
func generateV1Client(vns3Auth map[string]any, Logger Logger) (*cnv1.Client, error) {
	//check VNS3 controller host is specified
	host, hasHost := parseHost(vns3Auth)
	if !hasHost {
		return nil, fmt.Errorf("vns3 block requires host param")
	}
	//v1client uses url
	host = hostToUrl(host)
	//validate authentication
	username, hasUsername := parseUsername(vns3Auth)
	password, hasPassword := parsePassword(vns3Auth)
	token, hasToken := parseToken(vns3Auth)
	var v1Client *cnv1.Client
	var err error
	emptyString := ""
	if hasToken {
		Logger.Debug("using API Token auth for VNS3 v1 connection")
		if hasPassword && hasUsername {
			Logger.Warn("ignoring user and password in vns3 config")
		}
		v1Client, err = cnv1.NewClient(&emptyString, &emptyString, &token, &host)
	} else if hasPassword && hasUsername {
		Logger.Debug("using Basic auth for VNS3 v1 connection")
		v1Client, err = cnv1.NewClient(&username, &password, &emptyString, &host)
	} else {
		return nil, fmt.Errorf("vns3 config requires either username & password or token specified")
	}

	//check if any occurs parsing authentication
	if err != nil {
		return nil, err
	}
	return v1Client, nil
}

func generateVNS3Client(vns3Auth map[string]any, Logger Logger) (*cn.VNS3Client, error) {

	//check VNS3 controller host is specified
	host, hasHost := parseHost(vns3Auth)
	if !hasHost {
		return nil, fmt.Errorf("vns3 block requires host param")
	}

	//validate authentication
	username, hasUsername := parseUsername(vns3Auth)
	password, hasPassword := parsePassword(vns3Auth)
	apiToken, hasApiToken := parseToken(vns3Auth)

	var cfg *cn.Configuration
	if hasApiToken {
		Logger.Debug("using token auth for VNS3 connection")
		if hasPassword && hasUsername && hasApiToken {
			Logger.Warn("ignoring user and password in vns3 config")
		}
		cfg = cn.NewConfigurationWithAuth(host, cn.ContextAccessToken, apiToken)
	} else if hasPassword && hasUsername {
		Logger.Debug("using Basic auth for VNS3 connection")
		cfg = cn.NewConfigurationWithAuth(host, cn.ContextBasicAuth, cn.BasicAuth{
			UserName: username,
			Password: password,
		})
	} else {
		return nil, fmt.Errorf("vns3 config requires either username & password or api_token specified")
	}
	timeout, hasTimeout := vns3Auth["timeout"].(int)
	if !hasTimeout || timeout == 0 {
		timeout = 10
	}
	vns3Client := cn.NewVNS3Client(cfg, cn.ClientParams{
		Timeout: timeout,
		TLS:     false,
	})
	vns3Client.Log = Logger
	return vns3Client, nil
}

func hostToUrl(host string) string {
	if !strings.HasPrefix(host, "https") {
		return "https://" + host + ":8000/api"
	}
	return host
}

func parseHost(vns3Auth map[string]any) (string, bool) {
	host, has_host := vns3Auth["host"].(string)
	if !has_host || host == "" {
		return "", false
	} else {
		return host, true
	}
}

func parsePassword(vns3Auth map[string]any) (string, bool) {
	password, has_password := vns3Auth["password"].(string)
	if !has_password || password == "" {
		return "", false
	} else {
		return password, true
	}
}

func parseUsername(vns3Auth map[string]any) (string, bool) {
	vns3Username, hasUsername := vns3Auth["username"].(string)
	if !hasUsername {
		return "api", true
	} else if vns3Username == "" {
		return "api", true
	} else {
		return vns3Username, true
	}
}

func parseToken(vns3Auth map[string]any) (string, bool) {
	token, has_token := vns3Auth["token"].(string)
	if has_token && token != "" {
		return token, true
	}
	api_token, has_token := vns3Auth["api_token"].(string)
	if has_token && api_token != "" {
		return api_token, true
	}
	return "", false
}
