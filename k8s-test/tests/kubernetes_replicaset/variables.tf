resource "null_resource" "frontend_replicaset" {
  provisioner "local-exec" {
    command = "kubectl apply -f ${path.cwd}/frontend.yaml"
  }
}

resource "null_resource" "delay" {
  provisioner "local-exec" {
    command = "sleep 45"
  }
}


# Delay in order to get te resource creation complete
resource "null_resource" "get_pods" {
  depends_on = [
    null_resource.delay
  ]
  provisioner "local-exec" {
    command = "kubectl get replicasets"
  }
}
