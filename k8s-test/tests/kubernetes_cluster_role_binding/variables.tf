resource "null_resource" "create-cluster-role-binding" {
  provisioner "local-exec" {
    command = "kubectl apply -f ${path.cwd}/cluster_role.yaml"
  }
}

resource "null_resource" "delay" {
  provisioner "local-exec" {
    command = "sleep 45"
  }
}


# Delay in order to get te resource creation complete
resource "null_resource" "get_cluster_role_binding_jenkins" {
  depends_on = [
    null_resource.delay
  ]
  provisioner "local-exec" {
    command = "kubectl get clusterrolebindings jenkins"
  }
}
