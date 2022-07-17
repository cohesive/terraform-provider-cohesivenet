package cohesivenet

import (
	"context"
	"strconv"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFirewall() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFirewallRead,
		Schema: map[string]*schema.Schema{
			"response": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"private_ipaddress": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ipaddress": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_gateway": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"topology_checksum": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"vns3_version": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"topology_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"ntp_hosts": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"licensed": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"peered": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"asn": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"manager_id": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"overlay_ipaddress": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_token": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceFirewallRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cn.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	firewallResponse, err := c.GetFirewallRules()
	if err != nil {
		return diag.FromErr(err)
	}

	resp := make([]interface{}, 1, 1)
	row := make(map[string]interface{})

	row["private_ipaddress"] = firewallResponse.Firewall[0]
	row["public_ipaddress"] = firewallResponse.Firewall[1]

	resp[0] = row

	if err := d.Set("response", resp); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
