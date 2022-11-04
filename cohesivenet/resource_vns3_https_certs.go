package cohesivenet

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceHttpsCerts() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHttpsCertsCreate,
		ReadContext:   resourceHttpsCertsRead,
		UpdateContext: resourceHttpsCertsUpdate,
		DeleteContext: resourceHttpsCertsDelete,
		Schema: map[string]*schema.Schema{
			"cert_file": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File path to certificate",
			},
			"key_file": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File path to private key",
			},
			"vns3": &schema.Schema{
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: getVns3AuthSchema(),
				},
			},
		},
	}
}

func resourceHttpsCertsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	var diags diag.Diagnostics

	cert_file := d.Get("cert_file").(string)
	key_file := d.Get("key_file").(string)

	response, err := c.UpdateHttpsCerts(cert_file, key_file)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(response.Response.UUID)

	return diags

}

func resourceHttpsCertsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}

func resourceHttpsCertsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags

	//return resourceHttpsCertsRead(ctx, d, m)
}

func resourceHttpsCertsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}
