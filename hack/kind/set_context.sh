kubectl config --kubeconfig=kubeconfig set-cluster kind-hub --server=https://hub-control-plane:6443
kubectl config --kubeconfig=kubeconfig set-context kind-hub --cluster=kind-hub --user=kind-hub

kind load docker-image origin0119/cloud-chain:latest --name kind-hub