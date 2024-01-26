package cohesivenet

import (
	"context"
	"strconv"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlert() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlertCreate,
		ReadContext:   resourceAlertRead,
		UpdateContext: resourceAlertUpdate,
		DeleteContext: resourceAlertDelete,
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
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Alert Id",
			},
			"webhook_id": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Alert Id",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of deployed image",
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "URL of the image file to be imported",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "URL of a dockerfile that will be used to build the image",
			},
			"events": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "URL of a dockerfile that will be used to build the image",
			},
			"custom_properties": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Nested block for route attributes",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Alert events",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Alert events",
						},
					},
				},
			},
		},
	}
}

func resourceAlertCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}
	var diags diag.Diagnostics

	// synchronize creating an alert
	c.ReqLock.Lock()
	defer c.ReqLock.Unlock()

	name := d.Get("name").(string)
	url := d.Get("url").(string)
	webhookID := d.Get("webhook_id").(int)
	events := d.Get("events").([]interface{})
	customProperties := d.Get("custom_properties").([]interface{})
	enabled := d.Get("enabled").(bool)

	alert := cn.NewAlert{
		Name:             name,
		Url:              url,
		WebhookId:        webhookID,
		Events:           events,
		CustomProperties: customProperties,
		Enabled:          enabled,
	}

	newAlert, errCreateAlert := c.CreateAlert(alert)
	if errCreateAlert != nil {
		return diag.FromErr(errCreateAlert)
	}

	d.SetId(strconv.Itoa(newAlert.Response.ID))
	resourceAlertRead(ctx, d, m)

	return diags
}

/*
func resourceAlertRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}
*/

func resourceAlertRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}
	var diags diag.Diagnostics

	alertId := d.Id()

	_, errGetAlert := c.GetAlert(alertId)
	if errGetAlert != nil {
		return diag.FromErr(errGetAlert)
	}

	return diags
}

func resourceAlertUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}
	c.ReqLock.Lock()
	defer c.ReqLock.Unlock()

	if d.HasChange("name") ||
		d.HasChange("url") ||
		d.HasChange("webhook_id") ||
		d.HasChange("events") ||
		d.HasChange("custom_properties") ||
		d.HasChange("enabled") {

		name := d.Get("name").(string)
		url := d.Get("url").(string)
		webhookID := d.Get("webhook_id").(int)
		events := d.Get("events").([]interface{})
		customProperties := d.Get("custom_properties").([]interface{})
		enabled := d.Get("enabled").(bool)

		alert := cn.NewAlert{
			Name:             name,
			Url:              url,
			WebhookId:        webhookID,
			Events:           events,
			CustomProperties: customProperties,
			Enabled:          enabled,
		}
		alertId := d.Id()

		errUpdateAlert := c.UpdateAlert(alertId, alert)
		if errUpdateAlert != nil {
			return diag.FromErr(errUpdateAlert)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))

	}

	return resourceAlertRead(ctx, d, m)
}

/*
func resourceAlertDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}
*/

func resourceAlertDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}
	var diags diag.Diagnostics

	alertId := d.Id()

	errDeleteAlert := c.DeleteAlert(alertId)
	if errDeleteAlert != nil {
		return diag.FromErr(errDeleteAlert)
	}

	d.SetId("")

	return diags
}
