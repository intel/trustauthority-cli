GITCOMMIT := $(shell git describe --always)
BUILDDATE := $(shell TZ=UTC date +%Y-%m-%dT%H:%M:%S%z)
VERSION := "v0.3.0"
PROXY_EXISTS := $(shell if [[ "${https_proxy}" || "${http_proxy}" || "${no_proxy}" ]]; then echo 1; else echo 0; fi)
DOCKER_PROXY_FLAGS := ""
ifeq ($(PROXY_EXISTS),1)
    DOCKER_PROXY_FLAGS = --build-arg http_proxy="${http_proxy}" --build-arg https_proxy="${https_proxy}" --build-arg no_proxy="${no_proxy}"
else
    DOCKER_PROXY_FLAGS =
endif

tenantctl:
	mkdir -p out/
	env GOOS=linux GOSUMDB=off GOPROXY=direct go mod tidy && env GOOS=linux GOSUMDB=off GOPROXY=direct go build \
    -ldflags "-X intel/amber/tac/v1/utils.BuildDate=${BUILDDATE} -X intel/amber/tac/v1/utils.Version=${VERSION} -X intel/amber/tac/v1/utils.GitHash=${GITCOMMIT}" \
    -o out/tenantctl

installer: tenantctl
	mkdir -p out/installer
	cp dist/linux/install.sh out/installer/install.sh && chmod +x out/installer/install.sh
	cp out/tenantctl out/installer/tenantctl
	makeself out/installer out/tenantctl-$(VERSION).bin "Tenant CLI $(VERSION)" ./install.sh
	rm -rf installer

test:
	DOCKER_BUILDKIT=1 docker build ${DOCKER_PROXY_FLAGS} -f Dockerfile --target tester -t cli-unit-test:$(VERSION) .
	docker rmi cli-unit-test:$(VERSION)

all: clean test installer

clean:
	rm -rf out/*

.PHONY: installer all test clean
