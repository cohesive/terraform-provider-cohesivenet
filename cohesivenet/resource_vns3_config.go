package cohesivenet

// import (
// 	"context"
// 	"strconv"
// 	"strings"
// 	"time"
// 	"fmt"

// 	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet"
// 	macros "github.com/cohesive/cohesivenet-client-go/cohesivenet/macros"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// )


// func resourceVns3Config() *schema.Resource {
// 	return &schema.Resource{
// 		CreateContext: resourceConfigCreate,
// 		ReadContext:   resourceConfigRead,
// 		UpdateContext: resourceConfigUpdate,
// 		DeleteContext: resourceConfigDelete,
// 		Schema: map[string]*schema.Schema{
// 			"last_updated": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"host": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Optional: true,
// 			},
// 			"password": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Optional: true,
// 			},
// 			"token": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Optional: true,
// 			},
// 			"license_file": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Required: true,
// 			},
// 			"topology_name": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Required: true,
// 			},
// 			"controller_name": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Required: false,
// 			},
// 			"license_params": map[string]*schema.Schema{
// 				"default": &schema.Schema{
// 					Type:    schema.TypeBool,
// 					Default: true
// 				},
// 			},
// 			"keyset_params": map[string]*schema.Schema{
// 				"token": &schema.Schema{
// 					Type:     schema.TypeString,
// 					Required: true
// 				},
// 			},
// 			"peer_id": &schema.Schema{
// 				Type:     schema.TypeInt,
// 				Required: true,
// 			},
// 			"topology_checksum": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Computed: true,
// 			},
// 			"keyset_checksum": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Computed: true,
// 			},
// 			"licensed": &schema.Schema{
// 				Type:    schema.TypeBool,
// 				Computed: true,
// 			},
// 		},
// 	}
// }

// func resourceConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	var vns3 cohesivenet.VNS3Client
// 	// vns3 := m["vns3"].(cohesivenet.VNS3Client)

// 	// Warning or errors can be collected in a slice type
// 	var diags diag.Diagnostics

// 	host := d.Get("host").(string)
// 	password := d.Get("password").(string)
// 	token := d.Get("token").(string)

// 	// this is a lot of code to just determine vns3 client
// 	if host != "" || password != "" || token != "" {
// 		invalid := host == "" || (password == "" && token == "")
// 		if invalid {
// 			return diag.FromErr(fmt.Errorf("host and auth is required if host, password or token are passed"))
// 		}

// 		var cfg cohesivenet.Configuration
// 		if token != "" {
// 			cfg = cohesivenet.NewConfigurationWithAuth(host, cohesivenet.ContextAccessToken, token)
// 		} else {
// 			cfg = cohesivenet.NewConfigurationWithAuth(host, cohesivenet.ContextBasicAuth, cohesivenet.BasicAuth{
// 				UserName: "api",
// 				Password: password,
// 			})
// 		}

// 		vns3 = cohesivenet.NewVNS3Client(cfg, cohesivenet.ClientParams{
// 			Timeout: 3,
// 			TLS: false,
// 		})
// 		Logger := NewLogger(ctx)
// 		vns3.Log = Logger

// 	} else {
// 		vns3 = m["vns3"].(cohesivenet.VNS3Client)
// 	}

//     setupReq := macros.SetupRequest{
//         TopologyName: d.Get("topology_name").(string),
//         ControllerName: d.Get("controller_name").(string),
//         LicenseParams: cohesivenet.NewSetLicenseParametersRequest(true),
//         LicenseFile: "/Users/benplatta/code/cohesive/vns3-functional-testing/test-assets/license.txt",
//         PeerId: 1,
//         KeysetParams: cohesivenet.SetKeysetParamsRequest{
//             Token: "token",
//         },
//         WaitTimeout: 60*5,
//         KeysetTimeout: 60*5,
//     }

//     configDetail, setupErr := macros.SetupController(vns3, setupReq)

//     if setupErr != nil {
// 		Logger.Info(fmt.Sprintf("VNS3 Setup error: %+v", setupErr))
//     } else {
//         c := *configDetail
//         d, _ := c.MarshalJSON()
//         log.Printf("Setup success: %v", string(d))
// 		Logger.Info("VNS3 Setup success")
//     }


// 	d.SetId(strconv.Itoa(newEndpoint.ID))

// 	resourceConfigRead(ctx, d, m)

// 	return diags
// }

// /*
// func resourceConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	// Warning or errors can be collected in a slice type
// 	var diags diag.Diagnostics

// 	return diags
// }

// */
// func resourceConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	c := m.(map[string]interface{})["clientv1"].(cn.Client)
// 	// Warning or errors can be collected in a slice type
// 	var diags diag.Diagnostics

// 	endpointId := d.Id()

// 	endpoint, err := c.GetEndpoint(endpointId)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}

// 	//newEndpoint := endpoints.Response.(map[string]interface{})
// 	flatEndpoint := flattenEndpointData(endpoint)

// 	if err := d.Set("endpoint", flatEndpoint); err != nil {
// 		return diag.FromErr(err)
// 	}

// 	// always run
// 	//d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
// 	d.SetId(strconv.Itoa(endpoint.Response.ID))

// 	return diags
// }

// /*
// func resourceConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	return resourceConfigRead(ctx, d, m)
// }
// */

// func resourceConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	c := m.(map[string]interface{})["clientv1"].(cn.Client)

// 	endpointId := d.Id()

// 	if d.HasChange("endpoint") {

// 		endp := d.Get("endpoint").([]interface{})[0]
// 		endpoint := endp.(map[string]interface{})

// 		ep := cn.Endpoint{
// 			Name:                    endpoint["name"].(string),
// 			Description:             endpoint["description"].(string),
// 			Ipaddress:               endpoint["ipaddress"].(string),
// 			Secret:                  endpoint["secret"].(string),
// 			Pfs:                     endpoint["pfs"].(bool),
// 			Ike_version:             endpoint["ike_version"].(int),
// 			Nat_t_enabled:           endpoint["nat_t_enabled"].(bool),
// 			Extra_config:            endpoint["extra_config"].(string),
// 			Vpn_type:                endpoint["vpn_type"].(string),
// 			Route_based_int_address: endpoint["route_based_int_address"].(string),
// 			Route_based_local:       endpoint["route_based_local"].(string),
// 			Route_based_remote:      endpoint["route_based_remote"].(string),
// 		}

// 		_, err := c.UpdateEndpoint(endpointId, &ep)
// 		if err != nil {
// 			return diag.FromErr(err)
// 		}

// 		d.Set("last_updated", time.Now().Format(time.RFC850))
// 	}

// 	return resourceConfigRead(ctx, d, m)
// }

// func resourceConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	c := m.(map[string]interface{})["clientv1"].(cn.Client)

// 	// Warning or errors can be collected in a slice type
// 	var diags diag.Diagnostics

// 	endpointId := d.Id()

// 	err := c.DeleteEndpoint(endpointId)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}

// 	d.SetId("")

// 	return diags
// }