resource "cohesivenet_vns3_plugin_image" "image" {
    name = "test-tf-plugin"
    image_url  = ""
    description = "test-tf-plugin-description"
    command = "/usr/bin/supervisord"
 }