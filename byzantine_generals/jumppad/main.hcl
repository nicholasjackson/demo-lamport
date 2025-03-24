resource "network" "cloud" {
  subnet = "10.5.0.0/16"
}

resource "k8s_cluster" "k3s" {
  network {
    id = resource.network.cloud.meta.id
  }
}

output "KUBECONFIG" {
  value = resource.k8s_cluster.k3s.kube_config.path
}