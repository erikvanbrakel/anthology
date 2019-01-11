SHELL:=/bin/bash

GOCMD=go
GOGET=$(GOCMD) get
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY=anthology
PLATFORMS=darwin linux windows
ARCHITECTURES=386 amd64
PKG_DIR=./pkg

all: test build
build:
	$(GOBUILD) -o $(BINARY) -v
test:
	$(GOTEST) -v ./api/...
	$(GOTEST) -v ./tests/...
clean:
	$(GOCLEAN)
	rm -f $(BINARY)
	rm -rf $(PKG_DIR)
run:
	$(GOBUILD) -o $(BINARY) -v ./...
	./$(BINARY)

build-all:
	$(GOGET) github.com/mitchellh/gox
	mkdir -p $(PKG_DIR)

	CGO_ENABLED=0 gox -os="linux darwin windows" -arch="386 arm amd64" -osarch="!darwin/386 !darwin/arm" -output="$(PKG_DIR)/$(BINARY)_{{.OS}}_{{.Arch}}/$(BINARY)" .

zip-all:
	@cd $(PKG_DIR); \
	for fn in $$(ls **/*) ; \
	do \
		ZIPNAME=`dirname $$fn` ;\
		echo "Creating $$ZIPNAME"; \
		zip -j $$ZIPNAME.zip $$fn ;\
		rm $$fn; \
	done

publish: clean build-all zip-all
