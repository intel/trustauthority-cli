# Copyright (C) 2022 Intel Corporation
# SPDX-License-Identifier: BSD-3-Clause
SHELL := /bin/bash
GITCOMMIT := $(shell git describe --always)
BUILDDATE := $(shell TZ=UTC date +%Y-%m-%dT%H:%M:%S%z)
VERSION := v0.6.0
PROXY_EXISTS := $(shell if [[ "${https_proxy}" || "${http_proxy}" || "${no_proxy}" ]]; then echo 1; else echo 0; fi)
DOCKER_PROXY_FLAGS := ""
ifeq ($(PROXY_EXISTS),1)
    DOCKER_PROXY_FLAGS = --build-arg http_proxy="${http_proxy}" --build-arg https_proxy="${https_proxy}" --build-arg no_proxy="${no_proxy}"
else
    DOCKER_PROXY_FLAGS =
endif

tenantctl:
	mkdir -p out/
	DOCKER_BUILDKIT=1 docker build \
		${DOCKER_PROXY_FLAGS} \
		-f Dockerfile \
		--target builder  \
		--build-arg VERSION=${VERSION} \
		--build-arg COMMIT=${GITCOMMIT} \
		--build-arg BUILDDATE=${BUILDDATE} \
		-t tenantcli-build:${VERSION} \
		.
		docker run --rm -v `pwd`/out:/tmp/ tenantcli-build:${VERSION}
		docker rmi -f tenantcli-build:${VERSION}

installer: tenantctl
	mkdir -p out/installer
	cp dist/linux/install.sh out/installer/install.sh && chmod +x out/installer/install.sh
	cp out/tenantctl out/installer/tenantctl
	makeself out/installer out/tenantctl-$(VERSION)-$(GITCOMMIT).bin "Tenant CLI $(VERSION)" ./install.sh
	rm -rf installer

#This target once will work in CI
push-artifact: installer
	curl -sSf --user "$(ARTIFACTORY_USERNAME):$(ARTIFACTORY_PASSWORD)" -X PUT -T ./out/tenantctl-$(VERSION)-$(GITCOMMIT).bin  $(ARTIFACTORY)/releases/tenant-cli/tenantctl-$(VERSION)-$(GITCOMMIT).bin

test:
	DOCKER_BUILDKIT=1 docker build ${DOCKER_PROXY_FLAGS} -f Dockerfile --target tester -t cli-unit-test:$(VERSION) .

go-fmt: test
	docker run -i --rm cli-unit-test:$(VERSION) env GOOS=linux GOSUMDB=off /usr/local/go/bin/gofmt -l .

test-coverage: test
	docker run -i ${DOCKER_RUN_PROXY_FLAGS} --rm cli-unit-test:$(VERSION) /bin/bash -c "/usr/local/go/bin/go test ./... -coverprofile=cover.out; /usr/local/go/bin/go tool cover -func cover.out"

all: clean test installer test-coverage

clean:
	rm -rf out/*

.PHONY: installer all test clean go-fmt test-coverage push-artifact
