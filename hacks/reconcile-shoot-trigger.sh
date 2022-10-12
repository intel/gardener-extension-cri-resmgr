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
echo "Usage example: SHOOT=local2 OPERATION=retry hacks/reconcile-extension-trigger.sh - SHOOT/OPERATION are optional!"
set -x
OPERATION=${OPERATION:-reconcile}
SHOOT=${SHOOT:-local}
kubectl --context kind-gardener-local -n garden-local annotate shoot ${SHOOT} "gardener.cloud/operation=${OPERATION}" --overwrite
