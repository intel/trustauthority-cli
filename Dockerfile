# Copyright (C) 2022 Intel Corporation
# SPDX-License-Identifier: BSD-3-Clause
ARG TAAS_VERSION=v0.0.0
ARG TAAS_COMMIT=fffffff

FROM golang:1.18 AS tester
ARG TAAS_VERSION
ARG TAAS_COMMIT
WORKDIR /app
COPY . .
RUN BUILDDATE=$(TZ=UTC date +%Y-%m-%dT%H:%M:%S%z); \
        env GOOS=linux GOSUMDB=off go test ./... -coverprofile=cover.out; \
        env GOOS=linux GOSUMDB=off go tool cover -func=cover.out
