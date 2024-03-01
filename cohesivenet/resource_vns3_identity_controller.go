package cohesivenet

import (
	"context"
	"strconv"
	"time"

	cn "github.com/cohesive/cohesivenet-client-go/cohesivenet/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIdentityController() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityControllerCreate,
		ReadContext:   resourceIdentityControllerRead,
		UpdateContext: resourceIdentityControllerUpdate,
		DeleteContext: resourceIdentityControllerDelete,
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
			"port": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The LDAP port",
			},
			"encrypt": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to encrypt LDAP communication",
			},
			"encrypt_ldaps": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to encrypt LDAPS communication",
			},
			"encrypt_auth": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to encrypt auth",
			},
			"encrypt_auth_key": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to encrypt auth key",
			},
			"encrypt_auth_cert": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to encrypt auth cert",
			},
			"encrypt_verify_ca": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to encrypt verify CA",
			},
			"encrypt_ca_cert": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to encrypt CA cert",
			},
			"binddn": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The bind DN",
			},
			"bindpw": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The bind password",
			},
			"encrypt_auth_cert_data": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The auth cert data",
			},
			"encrypt_auth_cert_filename": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The auth cert filename",
			},
			"encrypt_auth_key_data": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The auth key data",
			},
			"encrypt_auth_key_filename": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The auth key filename",
			},
			"encrypt_ca_cert_data": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The CA cert data",
			},
			"encrypt_ca_cert_filename": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The CA cert filename",
			},
			"user_base": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The user base DN",
			},
			"user_id_attribute": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The user ID attribute",
			},
			"user_list_filter": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The user list filter",
			},
			"group_base": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The group base DN",
			},
			"group_id_attribute": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The group ID attribute",
			},
			"group_list_filter": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The group list filter",
			},
			"group_member_attribute": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The group member attribute",
			},
			"group_member_attr_format": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The group member attribute",
			},
			"group_search_scope": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The group list filter",
			},
			"otp": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The group member attribute",
			},
			"otp_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The group member attribute",
			},
		},
	}
}

func resourceIdentityControllerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	port := d.Get("port").(int)
	encrypt := d.Get("encrypt").(bool)
	encryptLdaps := d.Get("encrypt_ldaps").(bool)
	encryptAuth := d.Get("encrypt_auth").(bool)
	encryptAuthKey := d.Get("encrypt_auth_key").(bool)
	encryptAuthCert := d.Get("encrypt_auth_cert").(bool)
	encryptVerifyCa := d.Get("encrypt_verify_ca").(bool)
	encryptCaCert := d.Get("encrypt_ca_cert").(bool)
	binddn := d.Get("binddn").(string)
	bindpw := d.Get("bindpw").(string)
	encryptAuthCertData := d.Get("encrypt_auth_cert_data").(string)
	encryptAuthCertFilename := d.Get("encrypt_auth_cert_filename").(string)
	encryptAuthKeyData := d.Get("encrypt_auth_key_data").(string)
	encryptAuthKeyFilename := d.Get("encrypt_auth_key_filename").(string)
	encryptCaCertData := d.Get("encrypt_ca_cert_data").(string)
	encryptCaCertFilename := d.Get("encrypt_ca_cert_filename").(string)
	userBase := d.Get("user_base").(string)
	userIDAttribute := d.Get("user_id_attribute").(string)
	userListFilter := d.Get("user_list_filter").(string)
	groupBase := d.Get("group_base").(string)
	groupIDAttribute := d.Get("group_id_attribute").(string)
	groupListFilter := d.Get("group_list_filter").(string)
	groupMemberAttribute := d.Get("group_member_attribute").(string)
	groupMemberAttrFormat := d.Get("group_member_attr_format").(string)
	groupSearchScope := d.Get("group_search_scope").(string)
	otp := d.Get("otp").(bool)
	otpURL := d.Get("otp_url").(string)

	identityController := cn.NewIdentityController{
		Provider:                provider,
		Identifier:              identifier,
		Enabled:                 enabled,
		Secret:                  secret,
		RedirectUri:             redirectUri,
		AuthorizationEndpoint:   authorizationEndpoint,
		TokenEndpoint:           tokenEndpoint,
		UserinfoEndpoint:        userinfoEndpoint,
		JwksUri:                 jwksUri,
		RedirectHostname:        redirectHostname,
		ProviderUrl:             providerUrl,
		Issuer:                  issuer,
		Port:                    port,
		Encrypt:                 encrypt,
		EncryptLdaps:            encryptLdaps,
		EncryptAuth:             encryptAuth,
		EncryptAuthKey:          encryptAuthKey,
		EncryptAuthCert:         encryptAuthCert,
		EncryptVerifyCa:         encryptVerifyCa,
		EncryptCaCert:           encryptCaCert,
		Binddn:                  binddn,
		Bindpw:                  bindpw,
		EncryptAuthCertData:     encryptAuthCertData,
		EncryptAuthCertFilename: encryptAuthCertFilename,
		EncryptAuthKeyData:      encryptAuthKeyData,
		EncryptAuthKeyFilename:  encryptAuthKeyFilename,
		EncryptCaCertData:       encryptCaCertData,
		EncryptCaCertFilename:   encryptCaCertFilename,
		UserBase:                userBase,
		UserIDAttribute:         userIDAttribute,
		UserListFilter:          userListFilter,
		GroupBase:               groupBase,
		GroupIDAttribute:        groupIDAttribute,
		GroupListFilter:         groupListFilter,
		GroupMemberAttribute:    groupMemberAttribute,
		GroupMemberAttrFormat:   groupMemberAttrFormat,
		GroupSearchScope:        groupSearchScope,
		Otp:                     otp,
		OtpURL:                  otpURL,
	}

	_, errCreateIdentityController := c.CreateIdentityController(identityController)
	if errCreateIdentityController != nil {
		return diag.FromErr(errCreateIdentityController)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	resourceIdentityControllerRead(ctx, d, m)

	return diags
}

func resourceIdentityControllerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}
	var diags diag.Diagnostics

	_, errGetIdentityController := c.GetIdentityController()
	if errGetIdentityController != nil {
		return diag.FromErr(errGetIdentityController)
	}

	return diags
}

func resourceIdentityControllerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		d.HasChange("port") ||
		d.HasChange("encrypt_ldaps") ||
		d.HasChange("encrypt_auth") ||
		d.HasChange("encrypt_auth_key") ||
		d.HasChange("encrypt_auth_cert") ||
		d.HasChange("encrypt_verify_ca") ||
		d.HasChange("encrypt_ca_cert") ||
		d.HasChange("binddn") ||
		d.HasChange("bindpw") ||
		d.HasChange("encrypt_auth_cert_data") ||
		d.HasChange("encrypt_auth_cert_filename") ||
		d.HasChange("encrypt_auth_key_data") ||
		d.HasChange("encrypt_auth_key_filename") ||
		d.HasChange("encrypt_ca_cert_data") ||
		d.HasChange("encrypt_ca_cert_filename") ||
		d.HasChange("user_base") ||
		d.HasChange("user_id_attribute") ||
		d.HasChange("user_list_filter") ||
		d.HasChange("group_base") ||
		d.HasChange("group_id_attribute") ||
		d.HasChange("group_list_filter") ||
		d.HasChange("group_member_attribute") ||
		d.HasChange("group_member_attr_format") ||
		d.HasChange("group_search_scope") ||
		d.HasChange("otp") ||
		d.HasChange("otp_url") {

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
		port := d.Get("port").(int)
		encrypt := d.Get("encrypt").(bool)
		encryptLdaps := d.Get("encrypt_ldaps").(bool)
		encryptAuth := d.Get("encrypt_auth").(bool)
		encryptAuthKey := d.Get("encrypt_auth_key").(bool)
		encryptAuthCert := d.Get("encrypt_auth_cert").(bool)
		encryptVerifyCa := d.Get("encrypt_verify_ca").(bool)
		encryptCaCert := d.Get("encrypt_ca_cert").(bool)
		binddn := d.Get("binddn").(string)
		bindpw := d.Get("bindpw").(string)
		encryptAuthCertData := d.Get("encrypt_auth_cert_data").(string)
		encryptAuthCertFilename := d.Get("encrypt_auth_cert_filename").(string)
		encryptAuthKeyData := d.Get("encrypt_auth_key_data").(string)
		encryptAuthKeyFilename := d.Get("encrypt_auth_key_filename").(string)
		encryptCaCertData := d.Get("encrypt_ca_cert_data").(string)
		encryptCaCertFilename := d.Get("encrypt_ca_cert_filename").(string)
		userBase := d.Get("user_base").(string)
		userIDAttribute := d.Get("user_id_attribute").(string)
		userListFilter := d.Get("user_list_filter").(string)
		groupBase := d.Get("group_base").(string)
		groupIDAttribute := d.Get("group_id_attribute").(string)
		groupListFilter := d.Get("group_list_filter").(string)
		groupMemberAttribute := d.Get("group_member_attribute").(string)
		groupMemberAttrFormat := d.Get("group_member_attr_format").(string)
		groupSearchScope := d.Get("group_search_scope").(string)
		otp := d.Get("otp").(bool)
		otpURL := d.Get("otp_url").(string)

		identityController := cn.NewIdentityController{
			Provider:                provider,
			Identifier:              identifier,
			Enabled:                 enabled,
			Secret:                  secret,
			RedirectUri:             redirectUri,
			AuthorizationEndpoint:   authorizationEndpoint,
			TokenEndpoint:           tokenEndpoint,
			UserinfoEndpoint:        userinfoEndpoint,
			JwksUri:                 jwksUri,
			RedirectHostname:        redirectHostname,
			ProviderUrl:             providerUrl,
			Issuer:                  issuer,
			Port:                    port,
			Encrypt:                 encrypt,
			EncryptLdaps:            encryptLdaps,
			EncryptAuth:             encryptAuth,
			EncryptAuthKey:          encryptAuthKey,
			EncryptAuthCert:         encryptAuthCert,
			EncryptVerifyCa:         encryptVerifyCa,
			EncryptCaCert:           encryptCaCert,
			Binddn:                  binddn,
			Bindpw:                  bindpw,
			EncryptAuthCertData:     encryptAuthCertData,
			EncryptAuthCertFilename: encryptAuthCertFilename,
			EncryptAuthKeyData:      encryptAuthKeyData,
			EncryptAuthKeyFilename:  encryptAuthKeyFilename,
			EncryptCaCertData:       encryptCaCertData,
			EncryptCaCertFilename:   encryptCaCertFilename,
			UserBase:                userBase,
			UserIDAttribute:         userIDAttribute,
			UserListFilter:          userListFilter,
			GroupBase:               groupBase,
			GroupIDAttribute:        groupIDAttribute,
			GroupListFilter:         groupListFilter,
			GroupMemberAttribute:    groupMemberAttribute,
			GroupMemberAttrFormat:   groupMemberAttrFormat,
			GroupSearchScope:        groupSearchScope,
			Otp:                     otp,
			OtpURL:                  otpURL,
		}

		_, errUpdateIdentityController := c.UpdateIdentityController(identityController)
		if errUpdateIdentityController != nil {
			return diag.FromErr(errUpdateIdentityController)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))

	}

	return resourceIdentityControllerRead(ctx, d, m)
}

/*
func resourceIdentityVpnDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}
*/

func resourceIdentityControllerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, error := getV1Client(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
	}
	var diags diag.Diagnostics

	errDeleteIdentityController := c.DeleteIdentityController()
	if errDeleteIdentityController != nil {
		return diag.FromErr(errDeleteIdentityController)
	}

	d.SetId("")

	return diags
}
