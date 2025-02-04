resource "cohesivenet_vns3_plugin_image" "image" {
    name = "test-tf-plugin"
    image_url  = ""
    description = "test-tf-plugin-description"
    command = "/usr/bin/supervisord"
    documentation_url = "https://docs.cohesive.net/docs/network-edge-plugins/donamestuff/"
    support_url = "https://support.cohesive.net"
    version = "0.0.1"
    tags = {
        log-access = "true",
        ssh = "true"
    }
    metadata = {
        property1 = "string",
        property2 = "string"
  }
 }