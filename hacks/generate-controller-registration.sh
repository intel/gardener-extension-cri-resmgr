# Copyright 2022 Intel Corporation. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

chart="$(tar -C charts -c gardener-extension-cri-resmgr | gzip -n | base64 | tr -d '\n')"
OUT=examples/ctrldeploy-ctrlreg.yaml

#FOR DEBUG
#rm -rf /tmp/extract_dir && mkdir -p /tmp/extract_dir/ ; echo $chart | base64 -d  | gunzip | tar -xv -C /tmp/extract_dir && find /tmp/extract_dir

cat <<EOT > "$OUT"
---
apiVersion: core.gardener.cloud/v1beta1
kind: ControllerDeployment
metadata:
  name: cri-resmgr-extension
type: helm
providerConfig:
  chart: $chart

  values:
    ### For development purposes - set it to 0 (if you want to register extension but use local process with "make start").
    replicaCount: 1            

    ### For development purposes (tag) or for production use with internal registry (registry)
    ### possible to overwrite default values from charts/images.yaml. Works best togehter with "make TAG=mydevbranch build-images push-images"
    ### "image" overwrites extension (for seed) and "images" overwrite installation/agent (for shoot) images.
    # image:
    #   repository: localhost:5001/gardener-extension-cri-resmgr
    #   tag: mydevbranch
    #   pullPolicy: Always
    # imageVectorOverwrite: |
    #   images:
    #   - name: gardener-extension-cri-resmgr-installation
    #     tag: pawelbranch
    #     repository: localhost:5001/gardener-extension-cri-resmgr-installation
    #   - name: gardener-extension-cri-resmgr-agent
    #     tag: pawelbranch
    #     repository: localhost:5001/gardener-extension-cri-resmgr-agent
    
    ### For production or testing
    ### Default values are taken from 
    ### override defaults values from: charts/gardener-extensions-cri-resmgr/values.yaml and (internal) charts/gardener-extensions-cri-installation/values.yaml
    ### name should be one of: fallback/default/node.NODE/group.GROUP/force 
    # configs:
    #   fallback: |
    #     ... body ... - will be mounted to installation daemonset then copied to host and passed as --fallback-config to cri-resource-manager systemd unit
    #   default: |
    #     ... body ... - will be watched by node agent to push to all nodes (overwriting builtin and fallback)
    #   force: |
    #     ... body ... - will be mounted by installation daemonset then copied to host and passed as --force-config (overwriting builtin/fallback and config provided by agent)
---
apiVersion: core.gardener.cloud/v1beta1
kind: ControllerRegistration
metadata:
  name: cri-resmgr-extension
spec:
  deployment:
    # For development purpose - deploy the extensions before even shoots are created (or enabled)
    policy: Always
    deploymentRefs:
    - name: cri-resmgr-extension
  resources:
  - kind: Extension
    type: cri-resmgr-extension
    globallyEnabled: false
    reconcileTimeout: "60s"
EOT

echo "Successfully generated ControllerRegistration and ControllerDeployment example."
