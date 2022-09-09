kubectl -n garden-local annotate shoot local "confirmation.gardener.cloud/deletion=true" --overwrite
kubectl -n garden-local delete shoot local --wait=false
