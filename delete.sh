kind delete clusters gardener-local

# mkdir -p ~/work/
# git clone https://github.com/gardener/gardener ~/work/gardener
# cd ~/work/gardener
# git checkout v1.56.0
# cd -

make -C ~/work/gardener kind-up
kubectl cluster-info --context kind-gardener-local --kubeconfig ~/work/gardener/example/gardener-local/kind/local/kubeconfig
cp ~/work/gardener/example/gardener-local/kind/local/kubeconfig ~/.kube/config
kubectl get nodes

sleep 5

make -C ~/work/gardener/ gardener-up
helm list -n garden -a
make build-images push-images
# ./hacks/generate-controller-registration.sh
# make build-images push-images
kubectl apply -f ./examples/ctrldeploy-ctrlreg.yaml
kubectl get controllerregistrations.core.gardener.cloud cri-resmgr-extension
kubectl get controllerdeployments.core.gardener.cloud cri-resmgr-extension
kubectl get controllerinstallation.core.gardener.cloud

sleep 5

kubectl apply -f examples/shoot.yaml
kubectl get shoots -n garden-local --watch -o wide

kubectl patch shoot local -n garden-local -p '{"spec":{"extensions": [ {"type": "cri-resmgr-extension", "disabled": false} ] } }'
kubectl -n garden-local get secret local.kubeconfig -o jsonpath={.data.kubeconfig} | base64 -d > /tmp/kubeconfig-shoot-local.yaml
k konfig import -s  /tmp/kubeconfig-shoot-local.yaml

:extension
