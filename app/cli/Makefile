# Makefile
OCI_IMAGE_PATH = asia-south1-docker.pkg.dev/vapusdata-beta/vapusdata-ecosystem/vapusctl
TAG = local_$(shell whoami)
# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOLINT = golangci-lint
GOSEC = gosec
BUFCMD = buf
GOFMT = gofmt

# Build target
BINARY_NAME = vapusctl

all: test build ## Run tests and build

protogen: ## Generate protobuf
	cd ../../apis && make sync

build: ## Build the binary
# CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -installsuffix -ldflags="-w -s" -o $(BINARY_FILE)
	$(GOBUILD) -o $(BINARY_NAME) -ldflags="-X 'github.com/vapusdata-ecosystem/vapusai/app/cli.version=$(VERSION)'" main.go

install:
	$(GOCMD) install 

releaseoci:
	docker build --no-cache -t $(OCI_IMAGE_PATH):$(TAG) .
	rm -f $(BINARY_NAME)
	docker push $(OCI_IMAGE_PATH):$(TAG) 

test: ## Run tests
	$(GOTEST) -v ./...

clean: ## Clean the binary
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run: ## Run the binary
	$(GOBUILD) -o $(BINARY_NAME) -v 
	$(GOCMD) install  $(BINARY_NAME)
	$(BINARY_NAME)

runmain: ## Run main.go
	$(GOCMD) run main.go

lint: ## Run golangci-lint
	$(GOLINT) run

fmt: ## Run gofmt
	$(GOFMT) -w .

sec: ## Run gosec
	$(GOSEC) ./...

.PHONY: all build test clean run

sync:
	cd ../ && make sync
goget:
	echo argument is $(lib)
	$(GOGET) $(lib)
	make sync

goinstall:
	$(GOCMD) install $(lib)

help:
	@echo ''
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[$$()% 0-9a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)