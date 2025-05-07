.PHONY: init
# init env
init: init-api-tools
	go install github.com/vektra/mockery/v2@v2.43.2
	# using binary release for atlas, since ent schema handler is not included
	# in the community version anymore https://github.com/ariga/atlas/issues/2388#issuecomment-1864287189
	curl -sSf https://atlasgo.sh | sh -s -- -y

# initialize API tooling
.PHONY: init-api-tools
init-api-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.33.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
	go install github.com/bufbuild/buf/cmd/buf@v1.34.0
	go install github.com/envoyproxy/protoc-gen-validate@v1.0.1

init-env:
	# install Go 1.23.4
	ifeq ($(shell uname -s),darwin)
		ifeq ($(shell uname -m),arm64)
			curl -LO https://golang.org/dl/go1.23.darwin-arm64.tar.gz
			tar -C /usr/local -xzf go1.23.darwin-arm64.tar.gz
			rm go1.23.darwin-arm64.tar.gz
		else ifeq ($(shell uname -m),x86_64)
			curl -LO https://golang.org/dl/go1.23.darwin-amd64.tar.gz
			tar -C /usr/local -xzf go1.23.darwin-amd64.tar.gz
			rm go1.23.darwin-amd64.tar.gz
	else ifeq ($(shell uname -s),linux)
		curl -LO https://golang.org/dl/go1.23.darwin-amd64.tar.gz
		tar -C /usr/local -xzf go1.23.darwin-amd64.tar.gz
		rm go1.23.darwin-amd64.tar.gz
	endif
	# install Python 3.11.4
	ifeq ($(shell uname -s),darwin)
		curl -LO https://www.python.org/ftp/python/3.11.4/python-3.11.4-macos11.pkg
		open python-3.11.4-macos11.pkg
	else ifeq ($(shell uname -s),linux)
		curl -LO https://www.python.org/ftp/python/3.11.4/Python-3.11.4.tgz
		tar -xzf Python-3.11.4.tgz
		cd Python-3.11.4 && ./configure --enable-optimizations && make && sudo make install
		rm -rf Python-3.11.4 Python-3.11.4.tgz
	endif
	# install pipx
	ifeq ($(shell uname -s),darwin)
		brew install pipx
	else ifeq ($(shell uname -s),linux)
		python3 -m pip install --user pipx
		python3 -m pipx ensurepath
	endif
	# install poetry
	pipx install poetry
