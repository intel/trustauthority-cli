# Copyright (C) 2022 Intel Corporation
# SPDX-License-Identifier: BSD-3-Clause
ARG TAAS_VERSION=v0.0.0
ARG TAAS_COMMIT=fffffff

FROM golang:1.20.2 AS buildbase
ARG TAAS_VERSION
ARG TAAS_COMMIT
WORKDIR /app
COPY . .


FROM buildbase AS tester
ARG TAAS_VERSION
ARG TAAS_COMMIT
WORKDIR /app
RUN env GOOS=linux GOSUMDB=off go test ./...

FROM buildbase AS builder
ARG TAAS_VERSION
ARG TAAS_COMMIT
WORKDIR /app
RUN export PATH=$PATH:/usr/local/go/bin/ && \
  env GOOS=linux GOSUMDB=off GOPROXY=direct go mod tidy && go mod download && \
	env GOOS=linux GOSUMDB=off GOPROXY=direct go build \
    -ldflags "-X intel/amber/tac/v1/utils.BuildDate=${BUILDDATE} -X intel/amber/tac/v1/utils.Version=${VERSION} -X intel/amber/tac/v1/utils.GitHash=${GITCOMMIT}" \
    -o tenantctl
CMD ["cp", "/app/tenantctl", "/tmp/"]
