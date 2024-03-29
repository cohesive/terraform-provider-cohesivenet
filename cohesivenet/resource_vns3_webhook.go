package cohesivenet

import (
	"context"
	net "net/url"
	"strconv"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceWebhook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWebhookCreate,
		ReadContext:   resourceWebhookRead,
		UpdateContext: resourceWebhookUpdate,
		DeleteContext: resourceWebhookDelete,
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
				Computed:    true,
				Description: "Webhook id after creation",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of webhook",
			},
			"validate_cert": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Validate SSL/TLS certificate",
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "URL of the webhook enpoint",
			},
			"body": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom integration payload",
			},
			"events": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of alert events to be triggered",
			},
			"custom_properties": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of custom properties",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Custom property name",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Custom property value",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Custom property description",
						},
					},
				},
			},
			"headers": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "HTTP Headers to be included in the request",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "HTTP header name",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "HTTP header value",
						},
					},
				},
			},
			"parameters": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "HTTP Parameters to be included in the request",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "HTTP parameter name",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "HTTP parameter name",
						},
					},
				},
			},
		},
	}
}

func resourceWebhookCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}
	var diags diag.Diagnostics

	// synchronize creating an alert
	c.ReqLock.Lock()
	defer c.ReqLock.Unlock()

	name := d.Get("name").(string)
	validate_cert := d.Get("validate_cert").(bool)
	body := d.Get("body").(string)
	url := d.Get("url").(string)
	events := d.Get("events").([]interface{})

	customProperties := d.Get("custom_properties").([]interface{})
	customProps := make([]struct {
		Name        string `json:"name,omitempty"`
		Value       string `json:"value,omitempty"`
		Description string `json:"description,omitempty"`
	}, len(customProperties))

	for i, cp := range customProperties {
		custprop := cp.(map[string]interface{})
		customProps[i].Name = custprop["name"].(string)
		customProps[i].Value = custprop["value"].(string)
		customProps[i].Description = custprop["description"].(string)
	}

	headers := d.Get("headers").([]interface{})
	headrs := make([]struct {
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
	}, len(headers))

	for i, hd := range headers {
		header := hd.(map[string]interface{})
		headrs[i].Name = header["name"].(string)
		headrs[i].Value = header["value"].(string)
	}

	parameters := d.Get("parameters").([]interface{})
	params := make([]struct {
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
	}, len(parameters))

	for i, pm := range headers {
		param := pm.(map[string]interface{})
		params[i].Name = param["name"].(string)
		params[i].Value = param["value"].(string)
	}

	webhook := cn.NewWebhook{
		Name:             name,
		ValidateCert:     validate_cert,
		Body:             net.QueryEscape(body),
		URL:              url,
		Events:           events,
		CustomProperties: customProps,
		Headers:          headrs,
		Parameters:       params,
	}

	newWebhook, errCreateWebhook := c.CreateWebhook(webhook)
	if errCreateWebhook != nil {
		return diag.FromErr(errCreateWebhook)
	}

	d.SetId(strconv.Itoa(newWebhook.Response.ID))
	resourceWebhookRead(ctx, d, m)

	return diags
}

/*
func resourceWebhookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}
*/

func resourceWebhookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}
	var diags diag.Diagnostics

	webhookId := d.Id()

	_, errGetWebhook := c.GetWebhook(webhookId)
	if errGetWebhook != nil {
		return diag.FromErr(errGetWebhook)
	}

	return diags
}

/*
func resourceWebhookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}
*/

func resourceWebhookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}
	c.ReqLock.Lock()
	defer c.ReqLock.Unlock()

	if d.HasChange("name") ||
		d.HasChange("url") ||
		d.HasChange("validate_cert") ||
		d.HasChange("body") ||
		d.HasChange("headers") ||
		d.HasChange("parameters") ||
		d.HasChange("events") ||
		d.HasChange("custom_properties") {

		name := d.Get("name").(string)
		validate_cert := d.Get("validate_cert").(bool)
		body := d.Get("body").(string)
		url := d.Get("url").(string)
		events := d.Get("events").([]interface{})
		customProperties := d.Get("custom_properties")
		customProps := customProperties.([]struct {
			Name        string `json:"name,omitempty"`
			Value       string `json:"value,omitempty"`
			Description string `json:"description,omitempty"`
		})

		headers := d.Get("headers")
		headrs := headers.([]struct {
			Name  string `json:"name,omitempty"`
			Value string `json:"value,omitempty"`
		})

		parameters := d.Get("parameters")
		params := parameters.([]struct {
			Name  string `json:"name,omitempty"`
			Value string `json:"value,omitempty"`
		})

		webhook := cn.NewWebhook{
			Name:         name,
			ValidateCert: validate_cert,
			//Body:             body,
			Body:             net.QueryEscape(body),
			URL:              url,
			Events:           events,
			CustomProperties: customProps,
			Headers:          headrs,
			Parameters:       params,
		}
		webhookId := d.Id()

		errUpdateWebhook := c.UpdateWebhook(webhookId, webhook)
		if errUpdateWebhook != nil {
			return diag.FromErr(errUpdateWebhook)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))

	}

	return resourceWebhookRead(ctx, d, m)
}

/*
func resourceWebhookDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

		var diags diag.Diagnostics

		return diags
	}
*/
func resourceWebhookDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}
	var diags diag.Diagnostics

	webhookId := d.Id()

	errDeleteWebhook := c.DeleteWebhook(webhookId)
	if errDeleteWebhook != nil {
		return diag.FromErr(errDeleteWebhook)
	}

	d.SetId("")

	return diags
}
