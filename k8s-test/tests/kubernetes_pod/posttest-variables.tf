locals {
  filepath = "${path.cwd}/naked-pod.yml"
}

resource "null_resource" "named_test_resource" {
  provisioner "local-exec" {
    command = "kubectl delete -f ${local.filepath}"
  }
}

