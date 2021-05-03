resource "null_resource" "delete-frontend-replicaset" {
  provisioner "local-exec" {
    command = "kubectl delete -f ${path.cwd}/frontend.yaml"
  }
}

