GO          := go
GOBUILDOPTS := -v
BINPATH     := bin

build: build-converter build-spatialize build-tiler

build-converter:
	$(GO) build $(GOBUILDOPTS) -o "$(BINPATH)/grandine-converter" cmd/converter/*.go

build-spatialize:
	$(GO) build $(GOBUILDOPTS) -o "$(BINPATH)/grandine-spatialize" cmd/spatialize/*.go

build-tiler:
	$(GO) build $(GOBUILDOPTS) -o "$(BINPATH)/grandine-tiler" cmd/tiler/*.go 

test:
	$(GO) test ./...
