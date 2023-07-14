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
FROM golang:1.19.4-alpine3.16 AS builder

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
FROM gcr.io/distroless/base AS gardener-extension-cri-resmgr

COPY charts/internal /charts/internal
COPY --from=builder /go/bin/gardener-extension-cri-resmgr /
ENTRYPOINT ["/gardener-extension-cri-resmgr"]


### agnet and installation joined
FROM debian:latest as gardener-extension-cri-resmgr-installation-and-agent
# Please keep this in sync with CRI_RM_VERSION from Makefile!
COPY --from=intel/cri-resmgr-agent:v0.8.3 /bin/* /bin/
COPY Makefile .
RUN apt update
RUN apt install -y make wget
RUN make _install-binaries
ARG COMMIT=unset
ARG VERSION=unset
RUN bash -c "echo ${VERSION} >/VERSION"
RUN bash -c "echo ${COMMIT} >/COMMIT"
