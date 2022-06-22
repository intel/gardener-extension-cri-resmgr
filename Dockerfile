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
FROM golang:1.17.6-alpine AS builder

WORKDIR /go/src/github.com/intel/cri-resource-manager/packaging/gardener
COPY cmd .
COPY go.mod .
COPY go.sum .
RUN go install ./...

### extension
FROM alpine:3.15.0 AS gardener-extension-cri-rm

COPY charts /charts
COPY --from=builder /go/bin/gardener-extension-cri-rm /gardener-extension-cri-rm
ENTRYPOINT ["/gardener-extension-cri-rm"]

### installation
FROM ubuntu:22.04 AS gardener-extension-cri-rm-installation
RUN apt update -y && apt install -y make wget
COPY Makefile .
RUN make install-binaries
