package cohesivenet

import (
	"context"
	"strconv"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEndpoint() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEndpointsRead,
		Schema: map[string]*schema.Schema{
			"response": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"local_subnet": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"remote_subnet": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"endpointid": &schema.Schema{
										Type:     schema.TypeInt,
										Computed: true,
									},
									"endpoint_name": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"active": &schema.Schema{
										Type:     schema.TypeBool,
										Computed: true,
									},
									"enabled": &schema.Schema{
										Type:     schema.TypeBool,
										Computed: true,
									},
									"description": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"tunnel_params": &schema.Schema{
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"phase2": &schema.Schema{
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"connected": &schema.Schema{
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceEndpointsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//client := &http.Client{Timeout: 10 * time.Second}
	c := m.(*cn.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	//http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	//req, err := http.NewRequest("GET", fmt.Sprintf("%s/status/ipsec", "https://3.127.171.216:8000/api"), nil)
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//req.Header.Add("Api-Token", "771c844ecf0a2e0a9dd2c2a3071cfa7c1a06d7eed1f8664ce0995ec1b0824bee")
	//req.Header.Add("Api-Token", token)

	//r, err := client.Do(req)
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//defer r.Body.Close()
	//body, err := ioutil.ReadAll(r.Body)

	//endpoints := make([]map[string]interface{}, 0)
	//var endpoints []map[string]interface{}
	//errUnmarshal := json.Unmarshal([]byte(body), &endpoints)

	//err = json.NewDecoder(r.Body).Decode(&endpoints)

	//if errUnmarshal != nil {
	//	return diag.FromErr(err)
	//}

	//if err := d.Set("endpoints", endpoints); err != nil {
	//	return diag.FromErr(err)
	//}
	//endpointID := strconv.Itoa(d.Get("id").(int))
	endpointID := "1"
	c.GetEndpoint(endpointID)

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

/*
func dataSourceIpsecRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cn.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	orderID := strconv.Itoa(d.Get("id").(int))

	order, err := c.GetOrder(orderID)
	if err != nil {
		return diag.FromErr(err)
	}

	orderItems := flattenOrderItemsData(&order.Items)
	if err := d.Set("items", orderItems); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(orderID)

	return diags
}

func flattenOrderItemsData(orderItems *[]cn.OrderItem) []interface{} {
	if orderItems != nil {
		ois := make([]interface{}, len(*orderItems), len(*orderItems))

		for i, orderItem := range *orderItems {
			oi := make(map[string]interface{})

			oi["coffee_id"] = orderItem.Coffee.ID
			oi["coffee_name"] = orderItem.Coffee.Name
			oi["coffee_teaser"] = orderItem.Coffee.Teaser
			oi["coffee_description"] = orderItem.Coffee.Description
			oi["coffee_price"] = orderItem.Coffee.Price
			oi["coffee_image"] = orderItem.Coffee.Image
			oi["quantity"] = orderItem.Quantity

			ois[i] = oi
		}

		return ois
	}

	return make([]interface{}, 0)
}
*/
