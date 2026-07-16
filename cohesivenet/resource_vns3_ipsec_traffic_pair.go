package cohesivenet

import (
	"context"
	"strconv"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTrafficPair() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTrafficPairCreate,
		ReadContext:   resourceTrafficPairRead,
		UpdateContext: resourceTrafficPairUpdate,
		DeleteContext: resourceTrafficPairDelete,
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
				Description: "Endpoint ID to associate Trrafic Pair",
			},
			"remote_subnet": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Remote Subnet CIDR of Traffic Pair",
			},
			"local_subnet": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Local Subnet CIDR of Traffic Pair",
			},
			"ping_ipaddress": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IP Address to Send Keep Alive Pings",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Traffic Pair Description",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable / Disable Traffic Pair ",
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
			"ipsec_endpoint_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of IPsec Endpoint",
			},
		},
	}
}

func resourceTrafficPairCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	tp := cn.TrafficPair{
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

	trafficPairResponse, err := c.CreateTrafficPair(endpointId, remote_subnet, &tp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := trafficPairResponse.Response.ID

	d.Set("ipsec_endpoint_id", trafficPairResponse.Response.IpsecEndpointID)
	d.Set("remote_subnet", trafficPairResponse.Response.RemoteSubnet)
	d.Set("local_subnet", trafficPairResponse.Response.LocalSubnet)
	d.Set("ping_ipaddress", trafficPairResponse.Response.PingIpaddress)
	d.Set("ping_interval", trafficPairResponse.Response.PingInterval)
	d.Set("ping_interface", trafficPairResponse.Response.PingInterface)
	d.Set("description", trafficPairResponse.Response.Description)
	d.Set("enabled", trafficPairResponse.Response.Enabled)

	d.SetId(strconv.Itoa(id))

	resourceTrafficPairRead(ctx, d, m)

	return diags
}

func resourceTrafficPairRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	var diags diag.Diagnostics

	remoteSubnet := d.Get("remote_subnet").(string)
	endId, _ := d.GetOk("endpoint_id")
	endpointId := endId.(int)
	trafficPairId := d.Id()

	trafficPairResponse, err := c.GetTrafficPair(endpointId, remoteSubnet, trafficPairId)
	if err != nil {
		return diag.FromErr(err)
	}

	flatTrafficPair := flattenTrafficPairs(trafficPairResponse)

	// Handle empty response
	if len(flatTrafficPair) == 0 {
		d.SetId("")
		return diags
	}

	d.Set("remote_subnet", flatTrafficPair["remote_subnet"].(string))
	d.Set("local_subnet", flatTrafficPair["local_subnet"].(string))
	d.Set("ping_ipaddress", flatTrafficPair["ping_ipaddress"].(string))
	d.Set("ping_interval", flatTrafficPair["ping_interval"].(int))
	d.Set("ping_interface", flatTrafficPair["ping_interface"].(string))
	d.Set("description", flatTrafficPair["description"].(string))
	d.Set("enabled", flatTrafficPair["enabled"].(bool))

	return diags
}

func resourceTrafficPairUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	endId, _ := d.GetOk("endpoint_id")
	endpointId := endId.(int)

	trafficPairId := d.Id()

	if d.HasChange("remote_subnet") ||
		d.HasChange("local_subnet") ||
		d.HasChange("ping_ipaddress") ||
		d.HasChange("ping_interval") ||
		d.HasChange("ping_interface") ||
		d.HasChange("description") {

		remote_subnet := d.Get("remote_subnet").(string)
		local_subnet := d.Get("local_subnet").(string)
		ping_ipaddress := d.Get("ping_ipaddress").(string)
		ping_interval := d.Get("ping_interval").(int)
		ping_interface := d.Get("ping_interface").(string)
		enabled := d.Get("enabled").(bool)
		description := d.Get("description").(string)

		tp := cn.TrafficPair{
			Remote_Subnet:  remote_subnet,
			Local_Subnet:   local_subnet,
			Ping_Ipaddress: ping_ipaddress,
			Ping_Interval:  ping_interval,
			Ping_Interface: ping_interface,
			Enabled:        enabled,
			Description:    description,
		}

		_, err := c.UpdateTrafficPair(endpointId, trafficPairId, &tp)
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))

	}

	if d.HasChange("enabled") {

		enabled := d.Get("enabled").(bool)

		err := c.EnableDisableTrafficPair(endpointId, trafficPairId, enabled)
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceTrafficPairRead(ctx, d, m)
}

func resourceTrafficPairDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	var diags diag.Diagnostics

	endId, _ := d.GetOk("endpoint_id")
	endpointId := endId.(int)

	trafficPairId := d.Id()

	err := c.DeleteTrafficPair(endpointId, trafficPairId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func flattenTrafficPairs(newTrafficPair []cn.TrafficPair) map[string]interface{} {

	trafficPair := make(map[string]interface{}, len(newTrafficPair))
	i := 0

	for _, tp := range newTrafficPair {
		row := make(map[string]interface{})
		row["id"] = tp.ID
		row["local_subnet"] = tp.Local_Subnet
		row["remote_subnet"] = tp.Remote_Subnet
		row["ping_ipaddress"] = tp.Ping_Ipaddress
		row["ping_interval"] = tp.Ping_Interval
		row["ping_interface"] = tp.Ping_Interface
		row["enabled"] = tp.Enabled
		row["description"] = tp.Description
		trafficPair = row
		i++
	}

	return trafficPair
}
