# Makefile

# Go parameters
GOCMD = go
BUFCMD = buf

.PHONY: protogen sync api

protogen:
	rm -rf gen-go/protos
	rm -rf gen-go/vendor
	$(BUFCMD) dep update
	$(BUFCMD) lint
	$(BUFCMD) generate
	protoc-go-inject-tag -input gen-go/protos/models/v1alpha1/common.pb.go
	protoc-go-inject-tag -input gen-go/protos/models/v1alpha1/organization.pb.go
	protoc-go-inject-tag -input gen-go/protos/models/v1alpha1/users.pb.go
	protoc-go-inject-tag -input gen-go/protos/models/v1alpha1/vapusdata.pb.go
	protoc-go-inject-tag -input gen-go/protos/models/v1alpha1/vapusai.pb.go
	protoc-go-inject-tag -input gen-go/protos/vapusai-studio/v1alpha1/*.pb.go

sync: protogen
	cd gen-go && $(GOCMD) mod tidy && $(GOCMD) mod vendor
	

apisync: sync
	cd gen-go && $(GOCMD) run main.go