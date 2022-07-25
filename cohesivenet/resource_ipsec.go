package cohesivenet

import (
	"context"

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
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"endpoints": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_ip": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"secret": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"pfs": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ike_version": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"nat_t_enabled": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"extra_config": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpn_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_based_int_address": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_based_local": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_based_remote": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceEndpointsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
} /*
func resourceEndpointsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cn.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	endpoints := d.Get("endpoints").(interface{})

	newEndpoint := c.CreateEndpoints(endpoints)

	/*
		eps := []cn.EndpointResponse{}

		for _, endpoint := range endpoints {
			ep := endpoint.(map[string]interface{})

			co := i["coffee"].([]interface{})[0]
			coffee := co.(map[string]interface{})

			oi := cn.OrderItem{
				Coffee: cn.Coffee{
					ID: coffee["id"].(int),
				},
				Quantity: i["quantity"].(int),
			}

			ois = append(ois, oi)
		}

		o, err := c.CreateOrder(ois)
		if err != nil {
			return diag.FromErr(err)
		}


		d.SetId(strconv.Itoa(o.ID))


	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}
*/
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
