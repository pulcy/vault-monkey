PROJECT := vault-monkey
SCRIPTDIR := $(shell pwd)
ROOTDIR := $(shell cd $(SCRIPTDIR) && pwd)
VERSION:= $(shell cat $(ROOTDIR)/VERSION)
COMMIT := $(shell git rev-parse --short HEAD)

GOBUILDDIR := $(SCRIPTDIR)/.gobuild
SRCDIR := $(SCRIPTDIR)
BINDIR := $(ROOTDIR)/bin
VENDORDIR := $(SCRIPTDIR)/deps

ORGPATH := github.com/pulcy
ORGDIR := $(GOBUILDDIR)/src/$(ORGPATH)
REPONAME := $(PROJECT)
REPODIR := $(ORGDIR)/$(REPONAME)
REPOPATH := $(ORGPATH)/$(REPONAME)

MANIFESTTOOL := $(GOPATH)/bin/manifest-tool

# Magical rubbish to teach make what commas and spaces are.
EMPTY :=
SPACE := $(EMPTY) $(EMPTY)
COMMA := $(EMPTY),$(EMPTY)

LINUX_ARCH:=amd64 arm arm64 ppc64le s390x
PLATFORMS:=$(subst $(SPACE),$(COMMA),$(foreach arch,$(LINUX_ARCH),linux/$(arch)))

GOPATH := $(GOBUILDDIR)
GOVERSION := 1.10.0-alpine

ifndef GOOS
	GOOS := linux
endif
ifndef GOARCH
	GOARCH := amd64
endif

BINNAME := $(PROJECT)-$(GOOS)-$(GOARCH)
BIN := $(BINDIR)/$(BINNAME)

ifndef DOCKERIMAGE
	DOCKERIMAGE := $(PROJECT):dev
endif

SOURCES := $(shell find $(SRCDIR) -name '*.go')

.PHONY: all clean deps build

all: build-archs

build: $(BIN)

local:
	@${MAKE} -B GOOS=$(shell go env GOHOSTOS) GOARCH=$(shell go env GOHOSTARCH) build
	@ln -sf bin/$(PROJECT)-$(shell go env GOHOSTOS)-$(shell go env GOHOSTARCH) $(PROJECT)

build-archs:
	for arch in $(LINUX_ARCH); do \
		$(MAKE) -B DOCKERIMAGE=$(DOCKERIMAGE) GOARCH=$$arch docker ;\
	done

clean:
	rm -Rf bin $(GOBUILDDIR)

deps:
	@${MAKE} -B -s $(GOBUILDDIR)

$(GOBUILDDIR):
	@mkdir -p $(ORGDIR)
	@rm -f $(REPODIR) && ln -s ../../../.. $(REPODIR)
	@GOPATH=$(GOPATH) pulsar go flatten -V $(VENDORDIR)

update-vendor:
	@rm -Rf $(VENDORDIR)
	@pulsar go vendor -V $(VENDORDIR) \
		github.com/coreos/etcd/client \
		github.com/dchest/uniuri \
		github.com/dustin/go-humanize \
		github.com/giantswarm/retry-go \
		github.com/hashicorp/consul/api \
		github.com/hashicorp/go-rootcerts \
		github.com/hashicorp/hcl \
		github.com/hashicorp/vault/api \
		github.com/juju/errgo \
		github.com/kardianos/osext \
		github.com/kr/pretty \
		github.com/mitchellh/go-homedir \
		github.com/mitchellh/mapstructure \
		github.com/op/go-logging \
		github.com/ryanuber/columnize \
		github.com/spf13/cobra \
		github.com/spf13/pflag \
		github.com/YakLabs/k8s-client

$(BIN): $(GOBUILDDIR) $(SOURCES)
	@mkdir -p $(BINDIR)
	docker run \
		--rm \
		-v $(ROOTDIR):/usr/code \
		-e GOPATH=/usr/code/.gobuild \
		-e GOOS=$(GOOS) \
		-e GOARCH=$(GOARCH) \
		-e CGO_ENABLED=0 \
		-w /usr/code/ \
		golang:$(GOVERSION) \
		go build -installsuffix netgo -tags netgo -ldflags "-X main.projectVersion=$(VERSION) -X main.projectBuild=$(COMMIT)" -o /usr/code/bin/$(BINNAME) $(REPOPATH)

docker: $(BIN)
	docker build --build-arg arch=$(GOARCH) -t $(DOCKERIMAGE)-$(GOARCH) -f Dockerfile.build .

$(MANIFESTTOOL):
	go get github.com/estesp/manifest-tool

.PHONY: push-manifest
push-manifest: $(MANIFESTTOOL)
	echo Pushing image for platforms $(PLATFORMS)
	@$(MANIFESTTOOL) $(MANIFESTAUTH) push from-args \
    	--platforms $(PLATFORMS) \
    	--template $(DOCKERIMAGE)-ARCH \
    	--target $(DOCKERIMAGE)

.PHONY: release
release:
	@${MAKE} -B DOCKERIMAGE=pulcy/$(shell pulsar docker-tag) all push-manifest
