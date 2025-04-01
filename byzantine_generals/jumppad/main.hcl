resource "network" "cloud" {
  subnet = "10.5.0.0/16"
}

resource "k8s_cluster" "k3s" {
  network {
    id = resource.network.cloud.meta.id
  }
}

resource "ingress" "ui" {
  port = 8081

  target {
    resource = resource.k8s_cluster.k3s
    port     = 5173

    config = {
      service   = "lamport-ui"
      namespace = "default"
    }
  }
}

resource "ingress" "commander" {
  port = 8080

  target {
    resource = resource.k8s_cluster.k3s
    port     = 8080

    config = {
      service   = "lamport-commander"
      namespace = "default"
    }
  }
}

output "KUBECONFIG" {
  value = resource.k8s_cluster.k3s.kube_config.path
}