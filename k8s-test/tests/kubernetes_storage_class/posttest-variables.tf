resource "null_resource" "delete-storageclass" {
  provisioner "local-exec" {
    command = "kubectl delete -f ${path.cwd}/storageclass.yaml"
  }
}
