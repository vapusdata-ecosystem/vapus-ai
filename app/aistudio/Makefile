# Makefile

#input params
OCI_IMAGE_PATH = asia-south1-docker.pkg.dev/vapusdata-beta/vapusdata-ecosystem/aistudio
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

# Build target
BINARY_NAME = vapus_aistudio

all: test build

protogen:
	cd ../../apis && make sync

build-binary: protogen govendor
	cd cmd && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -a -installsuffix -ldflags="-w -s" -o  ../$(BINARY_NAME) 

build: build-binary
	docker build --no-cache -t $(OCI_IMAGE_PATH):$(TAG) .
	rm -f $(BINARY_NAME)

release: build
	docker push $(OCI_IMAGE_PATH):$(TAG) 

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

runmain:
	$(GOCMD) run cmd/main.go --conf=$(CONF) --debug=$(DEBUG)

lint:
	$(GOLINT) run --timeout=10m

sec:
	$(GOSEC) ./...

.PHONY: all build test clean run

prehook: lint sec test

sync: 
	cd ../ && make sync

govendor:
	cd ../ && make govendor

goget:
	echo argument is $(lib)
	$(GOGET) $(lib)
	make sync

goinstall:
	$(GOCMD) install $(lib)

