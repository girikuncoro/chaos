BINDIR := $(CURDIR)/bin
BINNAME ?= chaos

# go option
PKG := ./...
TAGS :=
LDFLAGS := -w -s
SRC  := $(shell find . -type f -name '*.go' -print)

.PHONY: all
all: build

.PHONY: build
build: $(BINDIR)/$(BINNAME)

$(BINDIR)/$(BINNAME): $(SRC)
	GO111MODULE=on go build $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o '$(BINDIR)'/$(BINNAME) ./cmd/chaos
