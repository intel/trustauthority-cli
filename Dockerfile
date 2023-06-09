# Copyright (C) 2022 Intel Corporation
# SPDX-License-Identifier: BSD-3-Clause

FROM golang:1.20.4 AS buildbase
ARG VERSION
ARG COMMIT
WORKDIR /app
COPY . .


FROM buildbase AS tester
ARG VERSION
ARG COMMIT
WORKDIR /app
RUN env GOOS=linux GOSUMDB=off go test ./...

FROM buildbase AS builder
ARG VERSION
ARG COMMIT
ARG BUILDDATE
WORKDIR /app
RUN export PATH=$PATH:/usr/local/go/bin/ && \
    env GOOS=linux GOSUMDB=off GOPROXY=direct go build \
    -ldflags "-X intel/amber/tac/v1/utils.BuildDate=${BUILDDATE} -X intel/amber/tac/v1/utils.Version=${VERSION} -X intel/amber/tac/v1/utils.GitHash=${COMMIT}" \
    -o tenantctl
CMD ["cp", "/app/tenantctl", "/tmp/"]