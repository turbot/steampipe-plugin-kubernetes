resource "null_resource" "delete-pdb" {
  provisioner "local-exec" {
    command = "kubectl delete -f ${path.cwd}/pdb.yaml"
  }
}

