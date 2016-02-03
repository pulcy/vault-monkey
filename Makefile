PROJECT := vault-monkey
SCRIPTDIR := $(shell pwd)
ROOTDIR := $(shell cd $(SCRIPTDIR) && pwd)
VERSION:= $(shell cat $(ROOTDIR)/VERSION)
COMMIT := $(shell git rev-parse --short HEAD)

GOBUILDDIR := $(SCRIPTDIR)/.gobuild
SRCDIR := $(SCRIPTDIR)
BINDIR := $(ROOTDIR)
VENDORDIR := $(SCRIPTDIR)/vendor

ORGPATH := git.pulcy.com/pulcy
ORGDIR := $(GOBUILDDIR)/src/$(ORGPATH)
REPONAME := $(PROJECT)
REPODIR := $(ORGDIR)/$(REPONAME)
REPOPATH := $(ORGPATH)/$(REPONAME)
BIN := $(BINDIR)/$(PROJECT)

GOPATH := $(GOBUILDDIR)
GOVERSION := 1.5.3

ifndef GOOS
	GOOS := linux
endif
ifndef GOARCH
	GOARCH := amd64
endif

SOURCES := $(shell find $(SRCDIR) -name '*.go')

.PHONY: all clean deps

all: $(BIN)

local:
	@${MAKE} -B GOOS=$(shell go env GOHOSTOS) GOARCH=$(shell go env GOHOSTARCH) $(BIN)

clean:
	rm -Rf $(BIN) $(GOBUILDDIR)

deps:
	@${MAKE} -B -s $(GOBUILDDIR)

$(GOBUILDDIR):
	@mkdir -p $(ORGDIR)
	@rm -f $(REPODIR) && ln -s ../../../.. $(REPODIR)

update-vendor:
	@rm -Rf $(VENDORDIR)
	@pulcy go vendor -V $(VENDORDIR) \
		github.com/dchest/uniuri \
		github.com/hashicorp/hcl \
		github.com/hashicorp/vault/api \
		github.com/juju/errgo \
		github.com/kardianos/osext \
		github.com/kr/pretty \
		github.com/mitchellh/mapstructure \
		github.com/op/go-logging \
		github.com/spf13/cobra \
		github.com/spf13/pflag

$(BIN): $(GOBUILDDIR) $(SOURCES)
	docker run \
	    --rm \
	    -v $(ROOTDIR):/usr/code \
	    -e GO15VENDOREXPERIMENT=1 \
	    -e GOPATH=/usr/code/.gobuild \
	    -e GOOS=$(GOOS) \
	    -e GOARCH=$(GOARCH) \
		-e CGO_ENABLED=0 \
	    -w /usr/code/ \
	    golang:$(GOVERSION) \
	    go build -a -installsuffix netgo -tags netgo -ldflags "-X main.projectVersion=$(VERSION) -X main.projectBuild=$(COMMIT)" -o /usr/code/$(PROJECT) $(REPOPATH)
