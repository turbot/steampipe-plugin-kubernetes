resource "null_resource" "delete-statefulset" {
  provisioner "local-exec" {
    command = "kubectl delete -f ${path.cwd}/statefulset.yaml"
  }
}
