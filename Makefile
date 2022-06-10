REGISTRY                    := v2.isvimgreg.com
EXTENSION_IMAGE_NAME        := gardener-extension-cri-rm
INSTALLATION_IMAGE_NAME     := gardener-extension-cri-rm-installation
VERSION                     := latest
CRI_RM_VERSION              := 0.6.1rc1
ARCHIVE_NAME                := cri-resource-manager-$(CRI_RM_VERSION).x86_64.tar.gz
CRI_RM_URL                  := https://github.com/intel/cri-resource-manager/releases/download/v$(CRI_RM_VERSION)/$(ARCHIVE_NAME)

.PHONY: start
start:
	go run ./cmd/gardener-extension-cri-rm --ignore-operation-annotation=true --leader-election=false

.PHONY: install
install:
	# TODO: Flags/version
	go install ./...

.PHONY: install-binaries
install-binaries:
	wget --directory-prefix=/cri-rm-installation https://github.com/intel/cri-resource-manager/releases/download/v$(CRI_RM_VERSION)/cri-resource-manager-$(CRI_RM_VERSION).x86_64.tar.gz
	tar -xvf /cri-rm-installation/cri-resource-manager-$(CRI_RM_VERSION).x86_64.tar.gz --directory /cri-rm-installation
	rm /cri-rm-installation/cri-resource-manager-$(CRI_RM_VERSION).x86_64.tar.gz
	

.PHONY: docker-images
docker-images:
	docker build -t $(REGISTRY)/$(EXTENSION_IMAGE_NAME):$(VERSION) -f Dockerfile --target $(EXTENSION_IMAGE_NAME) .
	docker build -t $(REGISTRY)/$(INSTALLATION_IMAGE_NAME):$(VERSION) -f Dockerfile --target $(INSTALLATION_IMAGE_NAME) .

.PHONY: publish-docker-images
publish-docker-images:
	docker push $(REGISTRY)/$(EXTENSION_IMAGE_NAME):$(VERSION)
	docker push $(REGISTRY)/$(INSTALLATION_IMAGE_NAME):$(VERSION)
