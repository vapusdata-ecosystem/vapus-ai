# Makefile

#input params
OCI_IMAGE_PATH = asia-south1-docker.pkg.dev/vapusdata-beta/vapusdata-ecosystem/webapp
TAG = local_$(shell whoami)

# Node.js parameters
NPMCMD = npm
NPMBUILD = $(NPMCMD) run build
NPMTEST = $(NPMCMD) test
NPMINSTALL = $(NPMCMD) install
NPMLINT = $(NPMCMD) run lint
NPMDEV = $(NPMCMD) run dev -- -p 9014

# Build target
BINARY_NAME = webapp

all: install build

install:
	$(NPMINSTALL)

build: install
	$(NPMBUILD)

build-docker: build
	docker build --no-cache -t $(OCI_IMAGE_PATH):$(TAG) .

release: build-docker
	docker push $(OCI_IMAGE_PATH):$(TAG)

test: install
	$(NPMTEST)

lint: install
	$(NPMLINT)

dev: install
	$(NPMDEV)

clean:
	rm -rf .next
	rm -rf node_modules
	rm -rf out

run: build
	$(NPMCMD) start

.PHONY: all build test clean run dev install build-docker release lint

prehook: lint test
