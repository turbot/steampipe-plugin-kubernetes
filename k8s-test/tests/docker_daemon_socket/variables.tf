# https://github.com/bridgecrewio/checkov/blob/master/checkov/kubernetes/checks/DockerSocketVolume.py
# https://github.com/bridgecrewio/checkov/blob/master/tests/kubernetes/checks/example_DockerSocketVolume/scope-2PASSED-1FAILED.yaml
# https://github.com/bridgecrewio/checkov/blob/master/tests/kubernetes/checks/example_DockerSocketVolume/cloudwatch-agent-1PASSED-1FAILED.yaml
# https://github.com/bridgecrewio/checkov/blob/master/tests/kubernetes/checks/test_DockerSocketVolume.py


#  name = "Do not expose the docker daemon socket to containers"
# Exposing the socket gives container information and increases risk of exploit
# read-only is not a solution but only makes it harder to exploit.
# Location: Pod.spec.volumes[].hostPath.path
# Location: CronJob.spec.jobTemplate.spec.template.spec.volumes[].hostPath.path
# Location: *.spec.template.spec.volumes[].hostPath.path
# supported_kind = ['Pod', 'Deployment', 'DaemonSet', 'StatefulSet', 'ReplicaSet', 'ReplicationController', 'Job', 'CronJob']


resource "null_resource" "namespace-amazon-cloudwatch" {
  provisioner "local-exec" {
    command = "kubectl create namespace amazon-cloudwatch"
  }
}

# deploy scope-2PASSED-1FAILED.yaml
resource "null_resource" "scope_failed" {
  provisioner "local-exec" {
    command = "kubectl apply -f ${path.cwd}/scope-2PASSED-1FAILED.yaml"
  }
}

# deploy cloudwatch-agent-1PASSED-1FAILED.yaml
resource "null_resource" "cloudwatch-agent" {
  depends_on = [
    null_resource.namespace-amazon-cloudwatch
  ]
  provisioner "local-exec" {
    command = "kubectl apply -f ${path.cwd}/cloudwatch-agent-1PASSED-1FAILED.yaml"
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

