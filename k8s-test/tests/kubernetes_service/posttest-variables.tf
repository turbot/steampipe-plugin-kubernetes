resource "null_resource" "delete-service" {
  provisioner "local-exec" {
    command = "kubectl delete -f ${path.cwd}/service.yaml"
  }
}
