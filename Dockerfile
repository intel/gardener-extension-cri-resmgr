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
# https://hub.docker.com/_/golang
FROM golang:1.22.4-alpine3.20 AS builder

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
# use latest from https://console.cloud.google.com/gcr/images/distroless/GLOBAL/static
#FROM gcr.io/distroless/static
# sha256:262ae336f8e9291f8edc9a71a61d5d568466edc1ea4818752d4af3d230a7f9ef Created Jan 1, 1, 1:24:00 AM
FROM gcr.io/distroless/static@sha256:262ae336f8e9291f8edc9a71a61d5d568466edc1ea4818752d4af3d230a7f9ef AS gardener-extension-cri-resmgr

COPY charts/internal/balloons /charts/internal/balloons
COPY --from=builder /go/bin/gardener-extension-cri-resmgr /
ENTRYPOINT ["/gardener-extension-cri-resmgr"]
