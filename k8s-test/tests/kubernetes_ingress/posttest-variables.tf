resource "null_resource" "delete_minimal_ingress" {
  provisioner "local-exec" {
    command = "kubectl delete -f ${path.cwd}/minimal-ingress.yaml"
  }
}

resource "null_resource" "delete_ingress_wildcard_host" {
  provisioner "local-exec" {
    command = "kubectl delete -f ${path.cwd}/ingress-wildcard-host.yaml"
  }
}

resource "null_resource" "delete_ingress_backend" {
  provisioner "local-exec" {
    command = "kubectl delete -f ${path.cwd}/ingress-resource-backend.yaml"
  }
}





