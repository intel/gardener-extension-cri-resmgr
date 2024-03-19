# Update dependencies procedure

1. Update go language toolchain runtime (go directive will be based on that)

```sh
wget https://go.dev/dl/go1.22.1.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.1.linux-amd64.tar.gz
go version
```

Warning: gardener project doesn't work under 1.22+, please install 1.21.8 for E2E tests

```
wget https://go.dev/dl/go1.21.8.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.21.8.linux-amd64.tar.gz
go version
```


2. Recreate go modules

rm go.mod go.sum

go mod init github.com/intel/gardener-extension-cri-resmgr
go mod tidy

3. Update base images in Dockerfile

- builder image: 

`FROM golang:1.22.1-alpine3.19 AS builder`

based on [golang images](https://hub.docker.com/_/golang`).

- installation image: 

`FROM debian:12.5 as gardener-extension-cri-resmgr-installation-and-agent`

based on [debian images](https://hub.docker.com/_/debian).


4. Update cri-resource-manager version in:

* Dockerfile:

`COPY --from=intel/cri-resmgr-agent:v0.9.0 /bin/* /bin/'

* Makefile:

`CRI_RM_VERSION=0.9.0`
