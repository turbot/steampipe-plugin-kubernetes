resource "null_resource" "create-service" {
  provisioner "local-exec" {
    command = "kubectl apply -f ${path.cwd}/service.yaml"
  }
}

resource "null_resource" "delay" {
  provisioner "local-exec" {
    command = "sleep 45"
  }
}


# Delay in order to get te resource creation complete
resource "null_resource" "get_services" {
  depends_on = [
    null_resource.delay
  ]
  provisioner "local-exec" {
    command = "kubectl get services"
  }
}
