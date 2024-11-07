# Copyright (C) 2022 Intel Corporation
# SPDX-License-Identifier: BSD-3-Clause
SHELL := /bin/bash
GITCOMMIT := $(shell git describe --always)
BUILDDATE := $(shell TZ=UTC date +%Y-%m-%dT%H:%M:%S%z)
VERSION := v1.1.9
PROXY_EXISTS := $(shell if [[ "${https_proxy}" || "${http_proxy}" || "${no_proxy}" ]]; then echo 1; else echo 0; fi)

trustauthorityctl:
	mkdir -p out/
	 env GOOS=linux CGO_CPPFLAGS="-D_FORTIFY_SOURCE=2" go build -buildmode=pie \
        -ldflags "-X intel/tac/v1/utils.BuildDate=${BUILDDATE} -X intel/tac/v1/utils.Version=${VERSION} -X intel/tac/v1/utils.GitHash=${GITCOMMIT} -linkmode=external -s -extldflags '-Wl,-z,relro,-z,now'"\
        -o out/trustauthorityctl

installer: trustauthorityctl
	mkdir -p out/installer
	cp dist/linux/install.sh out/installer/install.sh && chmod +x out/installer/install.sh
	mv out/trustauthorityctl out/installer/trustauthorityctl
	makeself out/installer out/trustauthorityctl-$(VERSION)-$(GITCOMMIT).bin "Intel Trust Authority $(VERSION)" ./install.sh
	rm -rf out/installer

push-artifact: installer
	curl -sSf --user "$(ARTIFACTORY_USERNAME):$(ARTIFACTORY_PASSWORD)" -X PUT -T ./out/trustauthorityctl-$(VERSION)-$(GITCOMMIT).bin  $(ARTIFACTORY)/releases/trust-authority-cli/trustauthorityctl-$(VERSION)-$(GITCOMMIT).bin

go-fmt:
	gofmt -l .

test-coverage:
	go test ./... -coverprofile=cover.out; go tool cover -func cover.out


all: clean test installer test-coverage

clean:
	rm -rf out/*

.PHONY: installer all test clean go-fmt test-coverage push-artifact
