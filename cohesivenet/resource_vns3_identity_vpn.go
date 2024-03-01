package cohesivenet

import (
	"context"
	"strconv"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIdentityVpn() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityVpnCreate,
		ReadContext:   resourceIdentityVpnRead,
		UpdateContext: resourceIdentityVpnUpdate,
		DeleteContext: resourceIdentityVpnDelete,
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
				Description: "Resource Id",
			},
			"identity_provider": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identity provider",
			},
			"identifier": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identifier",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enabled",
			},
			"secret": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Secret",
			},
			"redirect_uri": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Redirect Uri",
			},
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Host",
			},
			"authorization_endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Authorization Endpoint",
			},
			"token_endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Token Endpoint",
			},
			"userinfo_endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Userinfo Endpoint",
			},
			"jwks_uri": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Jwks Uri",
			},
			"redirect_hostname": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Redirect Hostname",
			},
			"provider_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Provider Url",
			},
			"issuer": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Issuer",
			},
		},
	}
}

func resourceIdentityVpnCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}
	var diags diag.Diagnostics

	// synchronize creating an alert
	c.ReqLock.Lock()
	defer c.ReqLock.Unlock()

	provider := d.Get("identity_provider").(string)
	identifier := d.Get("identifier").(string)
	enabled := d.Get("enabled").(bool)
	secret := d.Get("secret").(string)
	redirectUri := d.Get("redirect_uri").(string)
	authorizationEndpoint := d.Get("authorization_endpoint").(string)
	tokenEndpoint := d.Get("token_endpoint").(string)
	userinfoEndpoint := d.Get("userinfo_endpoint").(string)
	jwksUri := d.Get("jwks_uri").(string)
	redirectHostname := d.Get("redirect_hostname").(string)
	providerUrl := d.Get("provider_url").(string)
	issuer := d.Get("issuer").(string)

	identityVpn := cn.NewIdentityVpn{
		Provider:              provider,
		Identifier:            identifier,
		Enabled:               enabled,
		Secret:                secret,
		RedirectUri:           redirectUri,
		AuthorizationEndpoint: authorizationEndpoint,
		TokenEndpoint:         tokenEndpoint,
		UserinfoEndpoint:      userinfoEndpoint,
		JwksUri:               jwksUri,
		RedirectHostname:      redirectHostname,
		ProviderUrl:           providerUrl,
		Issuer:                issuer,
	}

	_, errCreateIdentityVpn := c.CreateIdentityVpn(identityVpn)
	if errCreateIdentityVpn != nil {
		return diag.FromErr(errCreateIdentityVpn)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	resourceIdentityVpnRead(ctx, d, m)

	return diags
}

func resourceIdentityVpnRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}
	var diags diag.Diagnostics

	_, errGetIdentityVpn := c.GetIdentityVpn()
	if errGetIdentityVpn != nil {
		return diag.FromErr(errGetIdentityVpn)
	}

	return diags
}

func resourceIdentityVpnUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}
	c.ReqLock.Lock()
	defer c.ReqLock.Unlock()

	if d.HasChange("identity_provider") ||
		d.HasChange("identifier") ||
		d.HasChange("enabled") ||
		d.HasChange("secret") ||
		d.HasChange("redirect_uri") ||
		d.HasChange("authorization_endpoint") ||
		d.HasChange("token_endpoint") ||
		d.HasChange("userinfo_endpoint") ||
		d.HasChange("jwks_uri") ||
		d.HasChange("redirect_hostname") ||
		d.HasChange("redirect_hostname") ||
		d.HasChange("issuer") {

		provider := d.Get("identity_provider").(string)
		identifier := d.Get("identifier").(string)
		enabled := d.Get("enabled").(bool)
		secret := d.Get("secret").(string)
		redirectUri := d.Get("redirect_uri").(string)
		authorizationEndpoint := d.Get("authorization_endpoint").(string)
		tokenEndpoint := d.Get("token_endpoint").(string)
		userinfoEndpoint := d.Get("userinfo_endpoint").(string)
		jwksUri := d.Get("jwks_uri").(string)
		redirectHostname := d.Get("redirect_hostname").(string)
		providerUrl := d.Get("provider_url").(string)
		issuer := d.Get("issuer").(string)

		identityVpn := cn.NewIdentityVpn{
			Provider:              provider,
			Identifier:            identifier,
			Enabled:               enabled,
			Secret:                secret,
			RedirectUri:           redirectUri,
			AuthorizationEndpoint: authorizationEndpoint,
			TokenEndpoint:         tokenEndpoint,
			UserinfoEndpoint:      userinfoEndpoint,
			JwksUri:               jwksUri,
			RedirectHostname:      redirectHostname,
			ProviderUrl:           providerUrl,
			Issuer:                issuer,
		}

		_, errUpdateIdentityVpn := c.UpdateIdentityVpn(identityVpn)
		if errUpdateIdentityVpn != nil {
			return diag.FromErr(errUpdateIdentityVpn)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))

	}

	return resourceIdentityVpnRead(ctx, d, m)
}

/*
func resourceIdentityVpnDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}
*/

func resourceIdentityVpnDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}
	var diags diag.Diagnostics

	errDeleteWebhook := c.DeleteIdentityVpn()
	if errDeleteWebhook != nil {
		return diag.FromErr(errDeleteWebhook)
	}

	d.SetId("")

	return diags
}
