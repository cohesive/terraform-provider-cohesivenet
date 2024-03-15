package cohesivenet

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLicenseUpgrade() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLicenseUpgradeCreate,
		ReadContext:   resourceLicenseUpgradeRead,
		UpdateContext: resourceLicenseUpgradeUpdate,
		DeleteContext: resourceLicenseUpgradeDelete,
		Schema: map[string]*schema.Schema{
			"license_upgrade_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "License upgrade file path",
			},
			"clientpack_ips": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comma separated list of clientpack ips",
			},
			"manager_ips": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comma separated list of manager ips",
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

func resourceLicenseUpgradeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	var diags diag.Diagnostics

	licenseUpgradeKeyPath := d.Get("license_upgrade_key").(string)
	clientpackIps := d.Get("clientpack_ips").(string)
	managerIps := d.Get("manager_ips").(string)

	err := c.CreateLicenseUpgrade(licenseUpgradeKeyPath, clientpackIps, managerIps)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(time.Now().Format(time.RFC850))

	return diags

}

func resourceLicenseUpgradeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	var diags diag.Diagnostics

	controllerLicense, err := c.GetControllerLicense()
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Controller License: %+v", controllerLicense)
	return diags
}

func resourceLicenseUpgradeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}

	var diags diag.Diagnostics

	if d.HasChange("license_upgrade_key") {

		licenseUpgradeKeyPath := d.Get("license_upgrade_key").(string)
		clientpackIps := d.Get("clientpack_ips").(string)
		managerIps := d.Get("manager_ips").(string)

		err := c.UpdateLicenseUpgrade(licenseUpgradeKeyPath, clientpackIps, managerIps)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	d.Set("last_updated", time.Now().Format(time.RFC850))

	return diags
}

func resourceLicenseUpgradeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	//Technically you cannot delete a license upgrade once it is applied to the controller.
	//But you may need to delete the resource from terraform state.
	d.SetId("")

	return diags
}
