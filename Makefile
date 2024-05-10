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
.PHONY: build clean e2e-test test start _install-binaries _build-agent-image build-images push-images _build-extension-image _build-installation-image

REGISTRY                         	:= localhost:5001/
EXTENSION_IMAGE_NAME             	:= gardener-extension-cri-resmgr
INSTALLATION_IMAGE_NAME          	:= gardener-extension-cri-resmgr-installation-and-agent
TAG                              := latest

# Please keep it up to date with agent image in charts/images.yaml
CRI_RM_VERSION                   := 0.9.0
CRI_RM_ARCHIVE_NAME              := cri-resource-manager-$(CRI_RM_VERSION).x86_64.tar.gz
CRI_RM_URL_RELEASE               := https://github.com/intel/cri-resource-manager/releases/download/v$(CRI_RM_VERSION)/$(CRI_RM_ARCHIVE_NAME)

# make start options
IGNORE_OPERATION_ANNOTATION 	 := false
# overwrite it if you want "make start" to read "configs" ConfigMap from Kubernetes
EXTENSION_CONFIGMAP_NAMESPACE    := ""

# Those two are lazy, because sometimes there is not .git context available
COMMIT:=`git rev-parse HEAD`
DIRTY:=`git diff --quiet || echo '-dirty'`
VERSION:=`git tag | sort -V | tail -1`

_go_generate:
	rm -rf ./pkg/consts/charts
	go generate ./...

_build:
	echo "Building ${VERSION}-${COMMIT}${DIRTY}"
	CGO_ENABLED=0 go build -ldflags="-X github.com/intel/gardener-extension-cri-resmgr/pkg/consts.Commit=${COMMIT}${DIRTY} -X github.com/intel/gardener-extension-cri-resmgr/pkg/consts.Version=${VERSION}" -v ./cmd/gardener-extension-cri-resmgr

_build_tests:
	go test -c -v  ./test/e2e/cri-resmgr-extension/. -o gardener-extension-cri-resmgr.e2e-tests
	go test -c -v ./pkg/controller/lifecycle -o ./gardener-extension-cri-resmgr.actuator.test
	go test -c -v ./pkg/configs -o ./gardener-extension-cri-resmgr.configs.test

build: _go_generate _build _build_tests

test:
	go generate ./...
	# Those tests (renders charts, uses env to read files) change CWD during execution (required because rely on charts and fixtures).
	go test  -v ./pkg/...

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
	@echo 'KUBECONFIG should point to kind-local gardener cluster e.g KUBECONFIG=$HOME/.kube/config'
	@echo ""
	@echo "Note3:"
	@echo "ControllerRegistration and ControllerDeployment CRDs must be already deployed to cluster (with globallyEnabled=False)"
	@echo 
	@echo "Note4:"
	@echo "Following labels are available: enable, reenable, disable"
	# Note seed 1 is used to keep order from simples to more complex cases (TODO to be replaced with SERIAL)
	ginkgo run -v --progress --seed 1 --slow-spec-threshold 2h --timeout 2h ./test/e2e/cri-resmgr-extension

start:
	go run ./cmd/gardener-extension-cri-resmgr --ignore-operation-annotation=$(IGNORE_OPERATION_ANNOTATION) --leader-election=false

_install-binaries:
	# WARNING: this should be run in container
	wget --no-check-certificate --directory-prefix=/cri-resmgr-installation $(CRI_RM_URL_RELEASE)
	tar -xvf /cri-resmgr-installation/$(CRI_RM_ARCHIVE_NAME) --directory /cri-resmgr-installation
	rm /cri-resmgr-installation/$(CRI_RM_ARCHIVE_NAME)

clean-images: 
	docker image rm $(REGISTRY)$(EXTENSION_IMAGE_NAME):$(TAG)
	docker image rm $(REGISTRY)$(INSTALLATION_IMAGE_NAME):$(TAG)

regenerate-charts:
	rm -rf ./pkg/consts/charts
	go generate ./...

_build-extension-image:
	@echo "Building extension image: commit=${COMMIT}${DIRTY} version=${VERSION} target=$(REGISTRY)$(EXTENSION_IMAGE_NAME):$(TAG)"
	docker build --build-arg COMMIT=${COMMIT}${DIRTY} --build-arg VERSION=${VERSION} -t $(REGISTRY)$(EXTENSION_IMAGE_NAME):$(TAG) -f Dockerfile --target $(EXTENSION_IMAGE_NAME) .

install-licenses:
	# go install github.com/google/go-licenses@latest
	# to be called from Docker image
	install -D LICENSE $(DESTDIR)/licenses/$$cmd/LICENSE
	go-licenses save ./cmd/ --ignore github.com/intel/gardener-extesion-cri-resmgr --save_path licenses

gen-licenses-csv:
	# go install github.com/google/go-licenses@latest
	@echo "WARNING: go-licenses require installed golang runtime matching one from go.mod go directive version!"
	go-licenses report ./cmd/gardener-extension-cri-resmgr >hacks/scans/licenses.csv

download-licenses-csv:
	# go install github.com/google/go-licenses@latest
	@echo "WARNING: go-licenses require installed golang runtime matching one from go.mod go directive version!"
	go-licenses save --save_path licenses ./...

_build-installation-image:
	@echo "Building installation image: commit=${COMMIT}${DIRTY} version=${VERSION} target=$(REGISTRY)$(INSTALLATION_IMAGE_NAME):$(TAG)"
	docker build --build-arg COMMIT=${COMMIT}${DIRTY} --build-arg VERSION=${VERSION} -t $(REGISTRY)$(INSTALLATION_IMAGE_NAME):$(TAG) -f Dockerfile --target $(INSTALLATION_IMAGE_NAME) .

dist: build build-images

build-images: regenerate-charts _build-extension-image _build-installation-image
	echo "Building ${VERSION}-${COMMIT}${DIRTY} done."


push-images:
	docker push $(REGISTRY)$(EXTENSION_IMAGE_NAME):$(TAG)
	docker push $(REGISTRY)$(INSTALLATION_IMAGE_NAME):$(TAG)
	echo "Images ${VERSION}-${COMMIT}${DIRTY} pushed."

generate-mocks:
	# go install go.uber.org/mock/mockgen@latest
	mockgen -destination=mocks/actuator.go -package=mocks github.com/gardener/gardener/extensions/pkg/controller/extension Actuator
	mockgen -destination=mocks/client.go -package=mocks sigs.k8s.io/controller-runtime/pkg/client Client,Reader,Object
	mockgen -destination=mocks/manager.go -package=mocks sigs.k8s.io/controller-runtime/pkg/manager Manager

generate-coverage:
	go test -coverprofile=coverage.out ./pkg/...
	go tool cover -html=coverage.out -o coverage.html    

lint:
	golangci-lint run 
