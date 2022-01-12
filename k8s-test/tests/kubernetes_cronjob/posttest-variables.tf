resource "null_resource" "delete-role" {
  provisioner "local-exec" {
    command = "kubectl delete -f ${path.cwd}/cronjob.yaml"
  }
}

