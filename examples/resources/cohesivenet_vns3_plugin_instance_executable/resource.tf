    resource  "cohesivenet_vns3_plugin_instance_executable" "suricata_update" {
    provider = cohesivenet.controller_1_updated
    instance_id = cohesivenet_vns3_plugin_instance.suricata_instance.id
    command = "Update"
    executable_path = "/opt/plugin-scripts/suricata-scripts.sh"
    timeout = 90
    run_count = 2

  }