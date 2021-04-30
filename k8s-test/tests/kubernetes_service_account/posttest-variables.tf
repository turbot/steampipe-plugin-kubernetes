resource "null_resource" "delete-service-account" {
  provisioner "local-exec" {
    command = "kubectl delete -f ${path.cwd}/service-account.yaml"
  }
}

