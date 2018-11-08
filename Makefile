GO          := go
GOBUILDOPTS ?= -v
GOTESTOPTS  := -v
BINPATH     := bin
CMDPREFIX   := github.com/thomersch/grandine/cmd

export CGO_CFLAGS=-I. -I/usr/local/include
export CGO_LDFLAGS=-L/usr/local/lib

build: build-converter build-inspect build-spatialize build-tiler

build-converter:
	$(GO) build $(GOBUILDOPTS) -o "$(BINPATH)/grandine-converter" $(CMDPREFIX)/converter

build-inspect:
	$(GO) build $(GOBUILDOPTS) -o "$(BINPATH)/grandine-inspect" $(CMDPREFIX)/inspect

build-spatialize:
	$(GO) build $(GOBUILDOPTS) -o "$(BINPATH)/grandine-spatialize" $(CMDPREFIX)/spatialize

build-tiler:
	$(GO) build $(GOBUILDOPTS) -o "$(BINPATH)/grandine-tiler" $(CMDPREFIX)/tiler

clean:
	rm '$(BINPATH)'/*

test:
	$(GO) test $(GOTESTOPTS) ./...

# retrieves deps for tests
test-deps:
	$(GO) get -t -u -v ./...
