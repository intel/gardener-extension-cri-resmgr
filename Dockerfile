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

### builder
FROM golang:1.21.5-alpine3.19 AS builder

WORKDIR /gardener-extension-cri-resmgr
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY cmd cmd
COPY cmd/gardener-extension-cri-resmgr/app cmd/gardener-extension-cri-resmgr/app
COPY pkg pkg
# only those two are required for building golang extension
COPY charts/images.go charts/images.go
COPY charts/images.yaml charts/images.yaml
ARG COMMIT=unset
ARG VERSION=unset
RUN CGO_ENABLED=0 go install -ldflags="-X github.com/intel/gardener-extension-cri-resmgr/pkg/consts.Commit=${COMMIT} -X github.com/intel/gardener-extension-cri-resmgr/pkg/consts.Version=${VERSION}" ./cmd/gardener-extension-cri-resmgr/...
# copying late saves time - no need to rebuild binary when only assest change
#COPY charts charts

### extension
# use latest from https://console.cloud.google.com/gcr/images/distroless/GLOBAL/base
#FROM gcr.io/distroless/base
FROM gcr.io/distroless/base@sha256:6c1e34e2f084fe6df17b8bceb1416f1e11af0fcdb1cef11ee4ac8ae127cb507c AS gardener-extension-cri-resmgr

COPY charts/internal /charts/internal
COPY --from=builder /go/bin/gardener-extension-cri-resmgr /
ENTRYPOINT ["/gardener-extension-cri-resmgr"]


### agnet and installation joined
FROM debian:12.2 as gardener-extension-cri-resmgr-installation-and-agent

WORKDIR /gardener-extension-cri-resmgr-installation-and-agent
# Please keep this in sync with CRI_RM_VERSION from Makefile!
COPY --from=intel/cri-resmgr-agent:v0.8.4 /bin/* /bin/
COPY Makefile .
RUN apt-get update && apt-get --no-install-recommends -y install make=4.3-4.1 wget=1.21.3-1+b2  && apt-get clean && rm -rf /var/lib/apt/lists/*
RUN make _install-binaries
ARG COMMIT=unset
ARG VERSION=unset
RUN bash -c "echo ${VERSION} >/VERSION" && bash -c "echo ${COMMIT} >/COMMIT"
