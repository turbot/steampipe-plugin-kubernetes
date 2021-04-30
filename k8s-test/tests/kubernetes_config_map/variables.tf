resource "null_resource" "create_config_map" {
  provisioner "local-exec" {
    command = "kubectl apply -f ${path.cwd}/config_map.yaml"
  }
}

resource "null_resource" "delay" {
  provisioner "local-exec" {
    command = "sleep 45"
  }
}


# Delay in order to get te resource creation complete
resource "null_resource" "get_config_map" {
  depends_on = [
    null_resource.delay
  ]
  provisioner "local-exec" {
    command = "kubectl get configmaps"
  }
}
