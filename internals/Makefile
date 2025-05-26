# Makefile

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOLINT = golangci-lint
GOSEC = gosec
GOVELNCHECK = govulncheck
GOVETCHECK = $(GOCMD) vet
GOFMT = gofmt


all: test build

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run: sync
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

runmain: gomodsync
	$(GOCMD) run main.go

lint:
	$(GOLINT) run

sec:
	$(GOSEC) ./...

.PHONY: all build test clean run

sync:
	cd ../apis && make sync
	$(GOCMD) mod tidy 

goget: gomodsync
	echo argument is $(lib)
	$(GOGET) $(lib)

goinstall:
	$(GOCMD) install $(lib)


vulncheck:
	$(GOVELNCHECK) ./...

gosec:
	$(GOSEC) ./...

govendor:
	cd ../apis && make sync
	$(GOCMD) mod tidy 
	$(GOCMD) mod vendor
