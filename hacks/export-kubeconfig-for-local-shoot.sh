echo 'source me!'
KUBECONFIG=~/.kube/config kubectl -n garden-local get secret local.kubeconfig -o jsonpath={.data.kubeconfig} | base64 -d > /tmp/kubeconfig-shoot-local.yaml
export KUBECONFIG=/tmp/kubeconfig-shoot-local.yaml
echo KUBECONFIG set to shoot local
