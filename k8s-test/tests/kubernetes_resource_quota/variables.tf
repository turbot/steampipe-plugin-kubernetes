resource "null_resource" "create_quota" {
  provisioner "local-exec" {
    command = "kubectl apply -f ${path.cwd}/resource_quota.yaml"
  }
}

resource "null_resource" "delay" {
  provisioner "local-exec" {
    command = "sleep 45"
  }
}


# Delay in order to get te resource creation complete
resource "null_resource" "get_quota" {
  depends_on = [
    null_resource.delay
  ]
  provisioner "local-exec" {
    command = "kubectl describe quota"
  }
}
