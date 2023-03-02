# Copyright (C) 2022 Intel Corporation
# SPDX-License-Identifier: BSD-3-Clause
ARG TAAS_VERSION=v0.0.0
ARG TAAS_COMMIT=fffffff

FROM golang:1.18 AS tester
ARG TAAS_VERSION
ARG TAAS_COMMIT
WORKDIR /app
COPY . .
RUN env GOOS=linux GOSUMDB=off go test ./...
