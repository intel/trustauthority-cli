# Copyright (C) 2022 Intel Corporation
# SPDX-License-Identifier: BSD-3-Clause
SHELL := /bin/bash
GITCOMMIT := $(shell git describe --always)
BUILDDATE := $(shell TZ=UTC date +%Y-%m-%dT%H:%M:%S%z)
VERSION := v0.6.0
PROXY_EXISTS := $(shell if [[ "${https_proxy}" || "${http_proxy}" || "${no_proxy}" ]]; then echo 1; else echo 0; fi)

tenantctl:
	mkdir -p out/
	 env GOOS=linux CGO_CPPFLAGS="-D_FORTIFY_SOURCE=2" go build -buildmode=pie \
        -ldflags "-X intel/amber/tac/v1/utils.BuildDate=${BUILDDATE} -X intel/amber/tac/v1/utils.Version=${VERSION} -X intel/amber/tac/v1/utils.GitHash=${GITCOMMIT} -linkmode=external -s -extldflags '-Wl,-z,relro,-z,now'"\
        -o out/tenantctl

installer: tenantctl
	mkdir -p out/installer
	cp dist/linux/install.sh out/installer/install.sh && chmod +x out/installer/install.sh
	cp out/tenantctl out/installer/tenantctl
	makeself out/installer out/tenantctl-$(VERSION)-$(GITCOMMIT).bin "Tenant CLI $(VERSION)" ./install.sh
	rm -rf installer

#This target once will work in CI
go-fmt:
	gofmt -l .

test-coverage:
	go test ./... -coverprofile=cover.out; go tool cover -func cover.out

all: clean test installer test-coverage

clean:
	rm -rf out/*

.PHONY: installer all test clean go-fmt test-coverage
