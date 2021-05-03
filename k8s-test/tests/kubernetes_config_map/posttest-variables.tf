resource "null_resource" "delete_config_map" {
  provisioner "local-exec" {
    command = "kubectl delete -f ${path.cwd}/config_map.yaml"
  }
}

