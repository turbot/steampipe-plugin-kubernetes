locals {
  filepath  = "${path.cwd}/naked-pod.yml"
  filepath1 = replace("${path.cwd}/naked-pod.yml", "terraform/test/", "")
}

# Create AWS > Lambda > Function
resource "local_file" "naked_pod" {
  filename = "${path.cwd}/naked-pod.yml"
  # sensitive_content = fil
}

output "filepath" {
  value = local.filepath
}
output "filepath1" {
  value = local.filepath1
}

resource "null_resource" "named_test_resource" {
  depends_on = [
    local_file.naked_pod
  ]
  provisioner "local-exec" {
    command = "kubectl create -f ${local.filepath1}"
  }
}

