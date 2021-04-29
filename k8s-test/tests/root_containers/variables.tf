# deploy rootContainersFAILED.yaml
resource "null_resource" "root_containers_failed" {
  provisioner "local-exec" {
    command = "kubectl create -f ${path.cwd}/rootContainersFAILED.yaml"
  }
}

# deploy rootContainersPASSED.yaml
resource "null_resource" "root_containers_passed" {
  provisioner "local-exec" {
    command = "kubectl create -f ${path.cwd}/rootContainersPASSED.yaml"
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


