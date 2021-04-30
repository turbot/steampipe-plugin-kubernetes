resource "null_resource" "create_role_binding" {
  provisioner "local-exec" {
    command = "kubectl apply -f ${path.cwd}/role_binding.yaml"
  }
}

resource "null_resource" "delay" {
  provisioner "local-exec" {
    command = "sleep 45"
  }
}


# Delay in order to get te resource creation complete
resource "null_resource" "get_role_binding" {
  depends_on = [
    null_resource.delay
  ]
  provisioner "local-exec" {
    command = "kubectl get rolebindings"
  }
}
