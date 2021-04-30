# https://github.com/bridgecrewio/checkov/blob/master/checkov/kubernetes/checks/HostPort.py
# https://github.com/bridgecrewio/checkov/tree/master/tests/kubernetes/checks/example_HostPort
# https://github.com/bridgecrewio/checkov/blob/master/tests/kubernetes/checks/test_HostPort.py

# delete DS-node-exporter-FAILED
resource "null_resource" "delete-ds-node-exporter-failed" {
  provisioner "local-exec" {
    command = "kubectl delete -f ${path.cwd}/DS-node-exporter-FAILED.yaml"
  }
}

# delete nginx-app-PASSED
resource "null_resource" "delete-nginx-app-passed" {
  provisioner "local-exec" {
    command = "kubectl delete -f ${path.cwd}/nginx-app-PASSED.yaml"
  }
}

resource "null_resource" "delete-namespace_monitoring" {
  depends_on = [
    null_resource.delete-ds-node-exporter-failed
  ]
  provisioner "local-exec" {
    command = "kubectl delete namespace monitoring"
  }
}




