set -x
kubectl patch shoot local -n garden-local -p '{"spec":{"extensions": [ {"type": "cri-resmgr-extension", "disabled": true} ] } }'
