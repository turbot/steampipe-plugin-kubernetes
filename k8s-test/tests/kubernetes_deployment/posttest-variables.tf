# https://github.com/bridgecrewio/checkov/blob/master/checkov/kubernetes/checks/HostPort.py
# https://github.com/bridgecrewio/checkov/tree/master/tests/kubernetes/checks/example_HostPort
# https://github.com/bridgecrewio/checkov/blob/master/tests/kubernetes/checks/test_HostPort.py

# delete DS-node-exporter-FAILED
resource "null_resource" "delete_deployment" {
  provisioner "local-exec" {
    command = "kubectl delete -f ${path.cwd}/nginx-deployment.yaml"
  }
}





