 resource  "vns3_plugin_instances" instance {
    name = "pluginname"
    plugin_id = vns3_plugin_images.image.id
    ip_address =  "198.51.100.11"
    description = "plugindescription"
    command = "/usr/bin/supervisord"
    environment = "HAENV_MODE=primary,HAENV_CLOUD=aws,HAENV_PEER_PUBLIC_IP=3.127.171.216,HAENV_SLEEP_TIME=15"

 }