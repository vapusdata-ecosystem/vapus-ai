# Makefile

#input params
OCI_IMAGE_PATH = asia-south1-docker.pkg.dev/vapusdata-beta/vapusai-studio/aitools
OCI_IMAGE_PATH_WINDOWS = asia-south1-docker.pkg.dev/vapusdata-beta/vapusai-studio/aitools_windows
TAG = local_$(shell whoami)

# Go parameters
BUFCMD = buf
POETRYCMD = poetry
BANDITCMD = bandit
PYLINTCMD = pylint
PYINSTALLER = pyinstaller

BINARY_NAME = aitools
BINARY_NAME_WINDOWS = aitools.exe

all: test build

protogen:
	cd ../../../apis && $(BUFCMD) generate

.PHONY: all build test clean run

build-binary: protogen 
	$(POETRYCMD) run $(PYINSTALLER) server.spec --clean --noconfirm

build: build-binary

	@if echo "$(shell powershell -Command "(systeminfo | Select-String '^OS Name').ToString().Split(':')[1].Trim()")" | grep -qi "Windows"; then \
		echo "Building for Windows..."; \
		docker build --no-cache -f Dockerfile_windows -t $(OCI_IMAGE_PATH_WINDOWS):$(TAG) .; \
		rm -f $(BINARY_NAME_WINDOWS)
	else \
		echo "Building for non-Windows..."; \
		docker build --no-cache -t $(OCI_IMAGE_PATH):$(TAG) .; \
		rm -f $(BINARY_NAME)
	fi

release: build
	docker push $(OCI_IMAGE_PATH):$(TAG) 

lint:
	$(POETRYCMD) run $(PYLINTCMD) controller/ models/ services/ utils/ core/ datastores/ helpers/ --rcfile=pylint.rc

sec:
	$(POETRYCMD) run $(BANDITCMD) -c bandit.toml -r -f json .

prehook: lint sec

pysync:
	$(POETRYCMD) install

runmain:
	poetry run python server.py --conf=$(CONF) 