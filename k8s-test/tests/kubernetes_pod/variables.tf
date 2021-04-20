locals {
  filepath = "${path.cwd}/naked-pod.yml"
}

output "filepath" {
  value = local.filepath
}

resource "null_resource" "named_test_resource" {
  provisioner "local-exec" {
    command = "kubectl create -f ${local.filepath}"
  }
}

