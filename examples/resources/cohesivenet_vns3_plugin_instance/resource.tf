 resource  "vns3_plugin_instance" instance {
    name = "pluginname"
    plugin_id = vns3_plugin_images.image.id
    ip_address =  "198.51.100.11"
    description = "plugindescription"
    command = "/usr/bin/supervisord"
    plugin_config = "/path/to/file"
 }