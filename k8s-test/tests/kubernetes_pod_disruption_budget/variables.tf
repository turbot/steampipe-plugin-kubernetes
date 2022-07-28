resource "null_resource" "create-pdb" {
  provisioner "local-exec" {
    command = "kubectl apply -f ${path.cwd}/pdb.yaml"
  }
}

resource "null_resource" "delay" {
  provisioner "local-exec" {
    command = "sleep 45"
  }
}

