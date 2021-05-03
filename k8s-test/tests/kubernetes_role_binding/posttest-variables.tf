resource "null_resource" "delete_role_binding" {
  provisioner "local-exec" {
    command = "kubectl delete -f ${path.cwd}/role_binding.yaml"
  }
}

