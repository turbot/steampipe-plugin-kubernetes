# delete rootContainersFAILED.yaml
resource "null_resource" "root_containers_failed" {
  provisioner "local-exec" {
    command = "kubectl delete -f ${path.cwd}/rootContainersFAILED.yaml"
  }
}

# delete rootContainersPASSED.yaml
resource "null_resource" "root_containers_passed" {
  provisioner "local-exec" {
    command = "kubectl delete -f ${path.cwd}/rootContainersPASSED.yaml"
  }
}



