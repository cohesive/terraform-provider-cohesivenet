package cohesivenet

import (
	"context"
	"strconv"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConfigRead,
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

func dataSourceConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(cn.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	response, err := c.GetConfig()
	if err != nil {
		return diag.FromErr(err)
	}

	resp := make([]interface{}, 1, 1)
	row := make(map[string]interface{})

	row["private_ipaddress"] = response.Config.PrivateIp
	row["public_ipaddress"] = response.Config.PublicIp
	row["subnet_gateway"] = response.Config.SubnetGateway
	row["topology_checksum"] = response.Config.TopologyChecksum
	row["vns3_version"] = response.Config.Vns3Version
	row["topology_name"] = response.Config.TopologyName
	row["ntp_hosts"] = response.Config.NtpHosts
	row["licensed"] = response.Config.Licensed
	row["peered"] = response.Config.Peered
	row["asn"] = response.Config.Asn
	row["manager_id"] = response.Config.ManagerId
	row["overlay_ipaddress"] = response.Config.OverlayIpaddress
	row["security_token"] = response.Config.SecurityToken
	resp[0] = row

	if err := d.Set("response", resp); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
