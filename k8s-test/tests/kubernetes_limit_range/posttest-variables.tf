resource "null_resource" "delete_limit" {
  provisioner "local-exec" {
    command = "kubectl delete -f ${path.cwd}/limit_range.yaml"
  }
}

