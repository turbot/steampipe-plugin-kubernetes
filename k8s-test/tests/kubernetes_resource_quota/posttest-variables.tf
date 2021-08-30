resource "null_resource" "delete_quota" {
  provisioner "local-exec" {
    command = "kubectl delete -f ${path.cwd}/quota.yaml"
  }
}

