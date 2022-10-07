package cohesivenet

import (
	"context"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet/v1"
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
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"key_file": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceHttpsCertsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(map[string]interface{})["clientv1"].(*cn.Client)

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
