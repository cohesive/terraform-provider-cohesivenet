resource "cohesivenet_vns3_identity_controller" "admin_oidc" {
  identity_provider = "oidc"
  identifier = "string"
  secret     = "string"
  enabled    = true
  redirect_hostname = "string"
  provider_url = "string"
  authorization_endpoint = "string"
  token_endpoint         = "string"
  userinfo_endpoint      = "string"
  jwks_uri               = "string"
  issuer = "string"
}