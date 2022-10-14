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
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Private IP of VNS3 instance",
						},
						"public_ipaddress": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public IP of VNS3 instance",
						},
						"subnet_gateway": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway IP of VNS3 subnet",
						},
						"topology_checksum": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Topology checksum of VNS3 network",
						},
						"vns3_version": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VNS3 controller version ",
						},
						"topology_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Topology name of VNS3 network",
						},
						"ntp_hosts": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP of time server",
						},
						"licensed": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Licensed flag",
						},
						"peered": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Peered flag",
						},
						"asn": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ASN of VNS3 controller",
						},
						"manager_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Id of VNS3 controller",
						},
						"overlay_ipaddress": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Overlay subnet of VNS3 network",
						},
						"security_token": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Token used to create VNS3 certficates",
						},
					},
				},
			},
		},
	}
}

func dataSourceConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

	var diags diag.Diagnostics

	response, err := c.GetConfig()
	if err != nil {
		return diag.FromErr(err)
	}

	resp := make([]interface{}, 1)
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
