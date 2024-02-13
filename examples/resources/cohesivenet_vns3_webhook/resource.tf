resource "cohesivenet_vns3_webhook" "webhook" {
  name = "custom webhook name"
  validate_cert = false
  url = "https://webhook url.com"
  body ="{\n    \"payload\": {\n      \"summary\": \"%%{default_message}\",\n      \"timestamp\": \"%%{event_date}\",\n      \"source\": \"%%{controller_name} @ %%{controller_public_ip}\",\n      \"severity\": \"critical\",\n      \"component\": \"VNS3\",\n      \"group\": \"EVENT\",\n      \"class\": \"vns3.%%{event_name}\",\n      \"custom_details\": {\n        \"Message\": \"%%{default_message}\"\n      }\n    }"
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
    description = "Description"
  }
    headers {
    name = "Name"
    value = "Value"

  }
    parameters {
    name = "Name"
    value = "Value"

  }
}