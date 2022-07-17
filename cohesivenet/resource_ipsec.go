package cohesivenet

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIpsec() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIpsecCreate,
		ReadContext:   resourceIpsecRead,
		UpdateContext: resourceIpsecUpdate,
		DeleteContext: resourceIpsecDelete,
		Schema: map[string]*schema.Schema{
			"ipsec": &schema.Schema{
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
					},
				},
			},
		},
	}
}

func resourceIpsecCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

/*
func resourceIpsecCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cn.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	ipsec := d.Get("ipsec").([]interface{})
	ois := []cn.OrderItem{}

	for _, ipsec := range ipsecs {
		i := ipsec.(map[string]interface{})

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

	return diags
}
*/

func resourceIpsecRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceIpsecUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceIpsecRead(ctx, d, m)
}

func resourceIpsecDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
