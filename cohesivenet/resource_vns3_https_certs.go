package cohesivenet

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceHttpsCerts() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHttpsCertsCreate,
		ReadContext:   resourceHttpsCertsRead,
		//There is no concept of read / update in API
		//UpdateContext: resourceHttpsCertsUpdate,
		DeleteContext: resourceHttpsCertsDelete,
		Schema: map[string]*schema.Schema{
			"cert_file": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Certificate file. Accepts absolute path",
			},
			"key_file": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,

				ForceNew:    true,
				Description: "Key file. Accepts Accepts absolute path",
			},
			"cert": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "Certificate file. Accepts file",
			},
			"key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "Key file. Accepts file",
			},
			"vns3": &schema.Schema{
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
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

	cert_file, hasCertFile := d.GetOk("cert_file")
	key_file, hasKeyFile := d.GetOk("key_file")
	cert, hasCert := d.GetOk("cert")
	key, hasKey := d.GetOk("key")

	if hasCertFile && hasKeyFile {
		response, err := c.UpdateHttpsCertsByFilepath(cert_file.(string), key_file.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(response.Response.UUID)
	} else if hasCert && hasKey {
		response, err := c.UpdateHttpsCertByValue(cert.(string), key.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(response.Response.UUID)
	} else {
		return diag.FromErr(fmt.Errorf("one of cert/key or cert_file/key_file must be specified"))
	}

	return diags

}

func resourceHttpsCertsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}

func resourceHttpsCertsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	//delete in the context of tf state
	d.SetId("")

	return diags
}
