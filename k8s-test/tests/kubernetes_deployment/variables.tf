
resource "null_resource" "create_deployment" {
  provisioner "local-exec" {
    command = "kubectl apply -f ${path.cwd}/nginx-deployment.yaml"
  }
}

resource "null_resource" "delay" {
  provisioner "local-exec" {
    command = "sleep 60"
  }
}


# Delay in order to get te resource creation complete
resource "null_resource" "get_deployments" {
  depends_on = [
    null_resource.delay
  ]
  provisioner "local-exec" {
    command = "kubectl get deployments"
  }
}


