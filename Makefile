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

REGISTRY                    := v2.isvimgreg.com/
EXTENSION_IMAGE_NAME        := gardener-extension-cri-resmgr
INSTALLATION_IMAGE_NAME     := gardener-extension-cri-resmgr-installation
VERSION                     := latest
CRI_RM_VERSION              := 0.6.1rc1
ARCHIVE_NAME                := cri-resource-manager-$(CRI_RM_VERSION).x86_64.tar.gz
CRI_RM_URL                  := https://github.com/intel/cri-resource-manager/releases/download/v$(CRI_RM_VERSION)/$(ARCHIVE_NAME)

.PHONY: build
build:
	go build -v ./cmd/gardener-extension-cri-resmgr
	go test -c -v ./test/integration/cri-resmgr-extension/...

.PHONY: test
test:
	go test -v ./test/integration/cri-resmgr-extension/...

.PHONY: start
start:
	go run ./cmd/gardener-extension-cri-resmgr --ignore-operation-annotation=true --leader-election=false

.PHONY: _install-binaries
_install-binaries:
	# WARNING: this should be run in container
	wget --directory-prefix=/cri-resmgr-installation https://github.com/intel/cri-resource-manager/releases/download/v$(CRI_RM_VERSION)/cri-resource-manager-$(CRI_RM_VERSION).x86_64.tar.gz
	tar -xvf /cri-resmgr-installation/cri-resource-manager-$(CRI_RM_VERSION).x86_64.tar.gz --directory /cri-resmgr-installation
	rm /cri-resmgr-installation/cri-resource-manager-$(CRI_RM_VERSION).x86_64.tar.gz
	
.PHONY: docker-images
docker-images:
	docker build -t $(REGISTRY)$(EXTENSION_IMAGE_NAME):$(VERSION) -f Dockerfile --target $(EXTENSION_IMAGE_NAME) .
	docker build -t $(REGISTRY)$(INSTALLATION_IMAGE_NAME):$(VERSION) -f Dockerfile --target $(INSTALLATION_IMAGE_NAME) .

.PHONY: publish-docker-images
publish-docker-images:
	docker push $(REGISTRY)$(EXTENSION_IMAGE_NAME):$(VERSION)
	docker push $(REGISTRY)$(INSTALLATION_IMAGE_NAME):$(VERSION)
