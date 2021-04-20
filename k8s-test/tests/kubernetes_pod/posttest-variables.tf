locals {
  filepath  = "${path.cwd}/naked-pod.yml"
  filepath1 = replace("${path.cwd}/naked-pod.yml", "terraform/test/", "")
}

resource "null_resource" "named_test_resource" {
  provisioner "local-exec" {
    command = "kubectl delete -f ${local.filepath1}"
  }
}

