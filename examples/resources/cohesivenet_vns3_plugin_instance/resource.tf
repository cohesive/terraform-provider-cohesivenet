resource  "cohesivenet_vns3_plugin_instance" instance {
   name = "pluginname"
   plugin_id = vns3_plugin_images.image.id
   ip_address =  "198.51.100.11"
   description = "plugindescription"
   command = "/usr/bin/supervisord"
   plugin_config = file("/path/to/file")

   plugin_config_files {
      filename = "Enable Rules"
      content = file("./enable.conf")
   }

   plugin_config_files {
      filename = "Disable Rules"
      content = file("./disable.conf")
   }
 }