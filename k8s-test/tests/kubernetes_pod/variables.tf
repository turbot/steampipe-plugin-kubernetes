# locals {
#   filepath = "${path.cwd}/naked-pod.yml"
# }

# output "filepath" {
#   value = local.filepath
# }

resource "null_resource" "naked-pod" {
  provisioner "local-exec" {
    command = "kubectl create -f ${path.cwd}/naked-pod.yml"
  }
}

resource "null_resource" "privileged-pod" {
  provisioner "local-exec" {
    command = "kubectl create -f ${path.cwd}/privileged-pod.yml"
  }
}

resource "null_resource" "pull-backoff" {
  provisioner "local-exec" {
    command = "kubectl create -f ${path.cwd}/pull-backoff.yml"
  }
}

