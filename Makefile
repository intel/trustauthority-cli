GITTAG := $(shell git describe --tags --abbrev=0 2> /dev/null)
GITCOMMIT := $(shell git describe --always)
BUILDDATE := $(shell TZ=UTC date +%Y-%m-%dT%H:%M:%S%z)
VERSION := "v0.1.0"

tenantctl:
	mkdir -p out/
	env GOOS=linux GOSUMDB=off GOPROXY=direct go mod tidy && env GOOS=linux GOSUMDB=off GOPROXY=direct go build -o out/tenantctl

installer: tenantctl
	mkdir -p out/installer
	cp dist/linux/install.sh out/installer/install.sh && chmod +x out/installer/install.sh
	cp out/tenantctl out/installer/tenantctl
	makeself out/installer out/tenantctl-$(VERSION).bin "Tenant CLI $(VERSION)" ./install.sh
	rm -rf installer

all: clean installer

clean:
	rm -rf out/*

.PHONY: installer all clean