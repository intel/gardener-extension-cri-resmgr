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
set -x
echo "Usage example: SHOOT=local2 hacks/disable-extension-in-local-shoot.sh"
SHOOT=${SHOOT:-local}
kubectl patch --context kind-gardener-local shoot ${SHOOT} -n garden-local -p '{"spec":{"extensions": [ {"type": "cri-resmgr-extension", "disabled": false} ] } }'
