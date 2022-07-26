package cohesivenet

import (
	"context"
	"strconv"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEndpoints() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEndpointsCreate,
		ReadContext:   resourceEndpointsRead,
		UpdateContext: resourceEndpointsUpdate,
		DeleteContext: resourceEndpointsDelete,
		Schema: map[string]*schema.Schema{
			"endpoint": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"ipaddress": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"secret": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"pfs": &schema.Schema{
							Type:     schema.TypeBool,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"ike_version": &schema.Schema{
							Type:     schema.TypeInt,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"nat_t_enabled": &schema.Schema{
							Type:     schema.TypeBool,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"extra_config": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"vpn_type": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"route_based_int_address": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"route_based_local": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
						"route_based_remote": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

/*
func resourceEndpointsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
*/

func resourceEndpointsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cn.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	endp := d.Get("endpoint").([]interface{})[0]
	endpoint := endp.(map[string]interface{})
	//ep := cn.Endpoint{}

	//	for _, end := range endpoints {
	//		e := end.(map[string]interface{})

	ep := cn.Endpoint{
		Name:                    endpoint["name"].(string),
		Description:             endpoint["description"].(string),
		Ipaddress:               endpoint["ipaddress"].(string),
		Secret:                  endpoint["secret"].(string),
		Pfs:                     endpoint["pfs"].(bool),
		Ike_version:             endpoint["ike_version"].(int),
		Nat_t_enabled:           endpoint["nat_t_enabled"].(bool),
		Extra_config:            endpoint["extra_config"].(string),
		Vpn_type:                endpoint["vpn_type"].(string),
		Route_based_int_address: endpoint["route_based_int_address"].(string),
		Route_based_local:       endpoint["route_based_local"].(string),
		Route_based_remote:      endpoint["route_based_remote"].(string),
	}

	//		ep = append(ep, endpoint)

	//	}

	_, err := c.CreateEndpoint(&ep)
	if err != nil {
		return diag.FromErr(err)
	}

	//d.SetId(strconv.Itoa(ep.Id))
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	//resourceEndpointsRead(ctx, d, m)

	return diags
}

func resourceEndpointsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceEndpointsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceEndpointsRead(ctx, d, m)
}

func resourceEndpointsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
