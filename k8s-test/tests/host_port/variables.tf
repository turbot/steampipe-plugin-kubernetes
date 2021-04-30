# https://github.com/bridgecrewio/checkov/blob/master/checkov/kubernetes/checks/HostPort.py
# https://github.com/bridgecrewio/checkov/tree/master/tests/kubernetes/checks/example_HostPort
# https://github.com/bridgecrewio/checkov/blob/master/tests/kubernetes/checks/test_HostPort.py

# https://kubernetes.io/docs/concepts/configuration/overview/
# Don’t specify a hostPort for a Pod unless it is absolutely necessary.
# When you bind a Pod to a hostPort, it limits the number of places the
# Pod can be scheduled, because each <hostIP, hostPort, protocol> combination
# must be unique.


resource "null_resource" "namespace_monitoring" {
  provisioner "local-exec" {
    command = "kubectl create namespace monitoring"
  }
}

# deploy nginx-app-PASSED
resource "null_resource" "nginx-app-passed" {
  provisioner "local-exec" {
    command = "kubectl create -f ${path.cwd}/nginx-app-PASSED.yaml"
  }
}

# deploy DS-node-exporter-FAILED
resource "null_resource" "ds-node-exporter-failed" {
  depends_on = [
    null_resource.namespace_monitoring
  ]
  provisioner "local-exec" {
    command = "kubectl create -f ${path.cwd}/DS-node-exporter-FAILED.yaml"
  }
}

resource "null_resource" "delay" {
  provisioner "local-exec" {
    command = "sleep 60"
  }
}


# Delay in order to get te resource creation complete
resource "null_resource" "get_pods" {
  depends_on = [
    null_resource.delay
  ]
  provisioner "local-exec" {
    command = "kubectl get pods"
  }
}


