
resource "null_resource" "create_minimal_ingress" {
  provisioner "local-exec" {
    command = "kubectl apply -f ${path.cwd}/minimal-ingress.yaml"
  }
}

resource "null_resource" "create_ingress_wildcard_host" {
  provisioner "local-exec" {
    command = "kubectl apply -f ${path.cwd}/ingress-wildcard-host.yaml"
  }
}

resource "null_resource" "create_ingress_backend" {
  provisioner "local-exec" {
    command = "kubectl apply -f ${path.cwd}/ingress-resource-backend.yaml"
  }
}

resource "null_resource" "delay" {
  provisioner "local-exec" {
    command = "sleep 60"
  }
}


# Delay in order to get te resource creation complete
resource "null_resource" "get_deployments" {
  depends_on = [
    null_resource.delay
  ]
  provisioner "local-exec" {
    command = "kubectl get ingresses"
  }
}


