#
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
#
.PHONY: build clean e2e-test test start _install-binaries cri-agent-docker-image docker-images publish-docker-images

REGISTRY                    := v2.isvimgreg.com/
EXTENSION_IMAGE_NAME        := gardener-extension-cri-resmgr
INSTALLATION_IMAGE_NAME     := gardener-extension-cri-resmgr-installation
AGENT_IMAGE_NAME            := gardener-extension-cri-resmgr-agent
VERSION                     := latest
CRI_RM_VERSION              := 0.7.2
CRI_RM_ARCHIVE_NAME         := cri-resource-manager-$(CRI_RM_VERSION).x86_64.tar.gz
CRI_RM_URL_RELEASE          := https://github.com/intel/cri-resource-manager/releases/download/v$(CRI_RM_VERSION)/$(CRI_RM_ARCHIVE_NAME)
CRI_RM_SRC_ARCHIVE_NAME     := vendored-cri-resource-manager-$(CRI_RM_VERSION).tar.gz
CRI_RM_URL_SRC              := https://github.com/intel/cri-resource-manager/releases/download/v$(CRI_RM_VERSION)/$(CRI_RM_SRC_ARCHIVE_NAME)

build:
	go build -v ./cmd/gardener-extension-cri-resmgr
	go test -c -v -o gardener-extension-cri-resmgr.e2e-tests ./test/e2e/cri-resmgr-extension/...
	go test -c -v ./cmd/gardener-extension-cri-resmgr

test:
	# Note: manual build required because this tests try to render charts project accessilbe only from root directory
	go test -c -v ./cmd/gardener-extension-cri-resmgr
	./gardener-extension-cri-resmgr.test --ginkgo.vv -test.v --ginkgo.progress

clean:
	go clean -cache -modcache -testcache
	rm cri-resmgr-extension.test
	rm gardener-extension-cri-resmgr

e2e-test:
	@echo "Note1:"
	@echo "Make sure following hosts are defined in etc/hosts"
	@echo "127.0.0.1 api.e2e-default.local.external.local.gardener.cloud"
	@echo "127.0.0.1 api.e2e-default.local.internal.local.gardener.cloud"
	@echo ""
	@echo "Note2:"
	@echo "KUBECONFIG should point to kind-local gardener cluster"
	@echo ""
	@echo "Note3:"
	@echo "ControllerRegistration and ControllerDeployment CRDs must be already deployed to cluster"
	@echo 
	@echo "Note4:"
	@echo "Following labels are available: enable, reenable, disable"
	# Note seed 1 is used to keep order from simples to more complex cases (TODO to be replaced with SERIAL)
	ginkgo run -v --progress --seed 1 --slow-spec-threshold 2h --timeout 2h ./test/e2e/cri-resmgr-extension

start:
	go run ./cmd/gardener-extension-cri-resmgr --ignore-operation-annotation=true

_install-binaries:
	# WARNING: this should be run in container
	wget --directory-prefix=/cri-resmgr-installation $(CRI_RM_URL_RELEASE)
	tar -xvf /cri-resmgr-installation/$(CRI_RM_ARCHIVE_NAME) --directory /cri-resmgr-installation
	rm /cri-resmgr-installation/$(CRI_RM_ARCHIVE_NAME)

cri-agent-docker-image:
	-mkdir tmpbuild
	wget --directory-prefix=tmpbuild -nc $(CRI_RM_URL_SRC)
	tar -C tmpbuild -xzvf tmpbuild/$(CRI_RM_SRC_ARCHIVE_NAME)
	# use exiting Dockerfile from cri-resource-manager source code
	docker build -t $(REGISTRY)$(AGENT_IMAGE_NAME):$(VERSION) -f tmpbuild/cri-resource-manager-$(CRI_RM_VERSION)/cmd/cri-resmgr-agent/Dockerfile  tmpbuild/cri-resource-manager-$(CRI_RM_VERSION)
	
docker-images: cri-agent-docker-image
	docker build -t $(REGISTRY)$(EXTENSION_IMAGE_NAME):$(VERSION) -f Dockerfile --target $(EXTENSION_IMAGE_NAME) .
	docker build -t $(REGISTRY)$(INSTALLATION_IMAGE_NAME):$(VERSION) -f Dockerfile --target $(INSTALLATION_IMAGE_NAME) .

publish-docker-images:
	docker push $(REGISTRY)$(EXTENSION_IMAGE_NAME):$(VERSION)
	docker push $(REGISTRY)$(INSTALLATION_IMAGE_NAME):$(VERSION)
	docker push $(REGISTRY)$(AGENT_IMAGE_NAME):$(VERSION)
