resource "null_resource" "delete-cluster-role" {
  provisioner "local-exec" {
    command = "kubectl delete -f ${path.cwd}/cluster_role.yaml"
  }
}

