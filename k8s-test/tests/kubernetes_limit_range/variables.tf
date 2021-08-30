resource "null_resource" "create_limit" {
  provisioner "local-exec" {
    command = "kubectl apply -f ${path.cwd}/limit_range.yaml --namespace=default"
  }
}

resource "null_resource" "delay" {
  provisioner "local-exec" {
    command = "sleep 45"
  }
}

