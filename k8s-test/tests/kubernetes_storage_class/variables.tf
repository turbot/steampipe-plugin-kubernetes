resource "null_resource" "create-storageclass" {
  provisioner "local-exec" {
    command = "kubectl apply -f ${path.cwd}/storageclass.yaml"
  }
}

resource "null_resource" "delay" {
  provisioner "local-exec" {
    command = "sleep 45"
  }
}


# Delay in order to get te resource creation complete
resource "null_resource" "get_storageclass" {
  depends_on = [
    null_resource.delay
  ]
  provisioner "local-exec" {
    command = "kubectl get storageclass"
  }
}
