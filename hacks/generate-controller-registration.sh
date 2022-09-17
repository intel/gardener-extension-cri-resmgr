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
  ### for development purposes (if you want to starte with "make start"
  #values:
  #  replicaCount: 0
  #  image:
  #    repository: localhost:5001/gardener-extension-cri-resmgr
  #    tag: latest
  #    pullPolicy: Always
  #  configs:
  #    foo2: |
  #      foo-2-from-ctrldeploy
---
apiVersion: core.gardener.cloud/v1beta1
kind: ControllerRegistration
metadata:
  name: cri-resmgr-extension
spec:
  deployment:
    deploymentRefs:
    - name: cri-resmgr-extension
  resources:
  - kind: Extension
    type: cri-resmgr-extension
    globallyEnabled: false
    reconcileTimeout: "60s"
EOT

echo "Successfully generated ControllerRegistration and ControllerDeployment example."
