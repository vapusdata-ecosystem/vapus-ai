# Makefile

#input params
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
	cd ../apis && make sync

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

lint:
	$(GOLINT) run --timeout=10m

sec:
	$(GOSEC) ./...

.PHONY: all build test clean run

prehook: lint sec test

sync: 
	cd ../internals && make sync
	$(GOCMD) mod tidy 
	cd aitools && make pysync

goget:
	echo argument is $(lib)
	$(GOGET) $(lib)
	make sync

goinstall:
	$(GOCMD) install $(lib)

govendor:
	cd ../internals && make govendor
	$(GOCMD) mod tidy 
	$(GOCMD) mod vendor
