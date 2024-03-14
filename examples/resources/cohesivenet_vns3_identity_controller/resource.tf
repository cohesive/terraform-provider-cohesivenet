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

resource "cohesivenet_vns3_identity_controller" "admin_ldap" {
  provider = cohesivenet.controller_1
  identity_provider = "ldap"
  enabled = true
  host = "ldap.google.com"
  port = 389
  binddn = "Supportive"
  bindpw = "password"
  encrypt = false
  user_base = "ou=People,dc=google,dc=com"
  user_id_attribute = "uid"
	user_list_filter = "*"
  group_base = "ou=Groups,dc=google,dc=com"
	group_id_attribute = "cn"
	group_list_filter = "*"
  group_member_attribute = "member"
	group_member_attr_format = "UserDN"
	group_search_scope = "subtree"

}