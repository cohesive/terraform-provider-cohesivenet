
resource "cohesivenet_vns3_alert" "alert" {
  webhook_id = cohesivenet_vns3_webhook.webhook.id
  name       = "alert name"
  url        = "https://webhook_url.com"
  events = [
                "tunnel_up",
                "tunnel_down",
                "process_change",
                "user_password_change",
                "controller_reboot",
                "controller_reset_defaults",
                "system_general",
                "clientpack_up",
                "clientpack_down",
                "bgp_session_up",
                "bgp_session_down",
                "wg_tunnel_up",
                "wg_tunnel_down",
                "fwset_dns_update"
  ]
  custom_properties {
    name = "Name"
    value = "Value"
  }

  depends_on = [ cohesivenet_vns3_webhook.webhook ]
}