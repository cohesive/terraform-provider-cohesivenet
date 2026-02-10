package cohesivenet

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTunnel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTunnelCreate,
		ReadContext:   resourceTunnelRead,
		UpdateContext: resourceTunnelUpdate,
		DeleteContext: resourceTunnelDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vns3": &schema.Schema{
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: getVns3AuthSchema(),
				},
			},
			"endpoint_id": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Endpoint ID to associate Tunnel",
			},
			"remote_subnet": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Remote Subnet CIDR of Tunnel",
				ValidateFunc: func(v interface{}, k string) ([]string, []error) {
					value := v.(string)
					if value == "" {
						return nil, []error{fmt.Errorf("remote_subnet cannot be empty")}
					}
					return nil, nil
				},
			},
			"local_subnet": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Local Subnet CIDR of Tunnel",
			},
			"ping_ipaddress": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IP Address to Send Keep Alive Pings",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Tunnel Description",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enabled / Disable Tunnel",
			},
			"ping_interval": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Interval between Keep Alive Pings",
			},
			"ping_interface": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Keep Alive Ping Interface (eth0/tun0)",
			},
			"tunnel_id": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of IPsec Tunnel",
			},
		},
	}
}

func resourceTunnelCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	var diags diag.Diagnostics

	remote_subnet := d.Get("remote_subnet").(string)
	local_subnet := d.Get("local_subnet").(string)
	ping_ipaddress := d.Get("ping_ipaddress").(string)
	ping_interval := d.Get("ping_interval").(int)
	ping_interface := d.Get("ping_interface").(string)
	enabled := d.Get("enabled").(bool)
	description := d.Get("description").(string)

	tun := cn.Tunnel{
		Remote_Subnet:  remote_subnet,
		Local_Subnet:   local_subnet,
		Ping_Ipaddress: ping_ipaddress,
		Ping_Interval:  ping_interval,
		Ping_Interface: ping_interface,
		Enabled:        enabled,
		Description:    description,
	}

	endId := d.Get("endpoint_id").(int)
	endpointId := strconv.Itoa(endId)

	tunnelResponse, err := c.CreateTunnel(endpointId, remote_subnet, &tun)
	if err != nil {
		return diag.FromErr(err)
	}

	flatTunnel := flattenTunnels(tunnelResponse)
	log.Printf("Create-flatTunnel%v", flatTunnel)

	id := flatTunnel["id"].(int)
	d.Set("tunnel_id", id)
	d.Set("remote_subnet", flatTunnel["remote_subnet"].(string))
	d.Set("local_subnet", flatTunnel["local_subnet"].(string))
	d.Set("ping_ipaddress", flatTunnel["ping_ipaddress"].(string))
	d.Set("ping_interval", flatTunnel["ping_interval"].(int))
	d.Set("ping_interface", flatTunnel["ping_interface"].(string))
	d.Set("description", flatTunnel["description"].(string))
	d.Set("enabled", flatTunnel["enabled"].(bool))

	d.SetId(strconv.Itoa(id))

	resourceTunnelRead(ctx, d, m)

	return diags
}

func resourceTunnelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	var diags diag.Diagnostics

	remoteSubnet := d.Get("remote_subnet").(string)
	if remoteSubnet == "" && d.HasChange("remote_subnet") {
		old, _ := d.GetChange("remote_subnet")
		remoteSubnet = old.(string)
	}

	endId, _ := d.GetOk("endpoint_id")
	endpointId := endId.(int)
	tunnelId := d.Id()

	tunnelResponse, err := c.GetTunnel(endpointId, remoteSubnet, tunnelId)
	if err != nil {
		log.Printf("[WARN] Failed to retrieve tunnel %s: %v", tunnelId, err)
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Ipsec Tunnel not found on controller. Tunnel will be recreated",
				Detail:   err.Error(),
			},
		}
	}

	if tunnelResponse == nil {
		log.Printf("[WARN] Ipsec Tunnel %s returned nil response", tunnelId)
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Ipsec Tunnel returned empty response",
				Detail:   fmt.Sprintf("Tunnel %s exists in state but controller returned no data", tunnelId),
			},
		}
	}

	flatTunnel := flattenTunnels(tunnelResponse)
	if len(flatTunnel) == 0 {
		// API didn't error but returned no matching tunnel
		log.Printf("[WARN] Tunnel %s not found", tunnelId)
		d.SetId("")
		return diags
	}

	d.Set("remote_subnet", flatTunnel["remote_subnet"].(string))
	d.Set("local_subnet", flatTunnel["local_subnet"].(string))
	d.Set("ping_ipaddress", flatTunnel["ping_ipaddress"].(string))
	d.Set("ping_interval", flatTunnel["ping_interval"].(int))
	d.Set("ping_interface", flatTunnel["ping_interface"].(string))
	d.Set("description", flatTunnel["description"].(string))
	d.Set("enabled", flatTunnel["enabled"].(bool))

	return diags
}

func resourceTunnelUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	endId, _ := d.GetOk("endpoint_id")
	endpointId := endId.(int)

	tunnelId := d.Id()

	if d.HasChange("remote_subnet") ||
		d.HasChange("local_subnet") ||
		d.HasChange("ping_ipaddress") ||
		d.HasChange("ping_interval") ||
		d.HasChange("ping_interface") ||
		d.HasChange("enabled") ||
		d.HasChange("description") {

		remote_subnet := d.Get("remote_subnet").(string)
		local_subnet := d.Get("local_subnet").(string)
		ping_ipaddress := d.Get("ping_ipaddress").(string)
		ping_interval := d.Get("ping_interval").(int)
		ping_interface := d.Get("ping_interface").(string)
		enabled := d.Get("enabled").(bool)
		description := d.Get("description").(string)

		tun := cn.Tunnel{
			Remote_Subnet:  remote_subnet,
			Local_Subnet:   local_subnet,
			Ping_Ipaddress: ping_ipaddress,
			Ping_Interval:  ping_interval,
			Ping_Interface: ping_interface,
			Enabled:        enabled,
			Description:    description,
		}

		_, err := c.UpdateTunnel(endpointId, tunnelId, remote_subnet, &tun)
		if err != nil {
			diags := diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Tunnel update failed, refreshing state from controller",
					Detail:   err.Error(),
				},
			}

			//Refresh state to match controller.
			readDiags := resourceTunnelRead(ctx, d, m)
			return append(diags, readDiags...)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceTunnelRead(ctx, d, m)
}

func resourceTunnelDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	var diags diag.Diagnostics

	endId, _ := d.GetOk("endpoint_id")
	endpointId := endId.(int)

	tunnelId := d.Id()

	err := c.DeleteTunnel(endpointId, tunnelId)
	if err != nil {
		// If tunnel/endpoint not found, consider it already deleted. Forces state to be rewritten.
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "404") {
			log.Printf("[WARN] Tunnel %s or endpoint %d not found, removing from state", tunnelId, endpointId)
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func flattenTunnels(newTunnel []cn.NewTunnel) map[string]interface{} {

	tunnel := make(map[string]interface{}, len(newTunnel))
	i := 0

	for _, tn := range newTunnel {
		row := make(map[string]interface{})
		row["id"] = tn.ID
		row["local_subnet"] = tn.LocalSubnet
		row["remote_subnet"] = tn.RemoteSubnet
		row["ping_ipaddress"] = tn.PingIpaddress
		row["ping_interval"] = tn.PingInterval
		row["ping_interface"] = tn.PingInterface
		row["enabled"] = tn.Enabled
		row["description"] = tn.Description
		tunnel = row
		i++
	}

	return tunnel
}
