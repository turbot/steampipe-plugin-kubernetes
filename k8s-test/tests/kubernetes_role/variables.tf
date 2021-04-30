resource "null_resource" "create-role" {
  provisioner "local-exec" {
    command = "kubectl apply -f ${path.cwd}/role.yaml"
  }
}

resource "null_resource" "delay" {
  provisioner "local-exec" {
    command = "sleep 45"
  }
}


# Delay in order to get te resource creation complete
resource "null_resource" "get_role" {
  depends_on = [
    null_resource.delay
  ]
  provisioner "local-exec" {
    command = "kubectl get roles"
  }
}
