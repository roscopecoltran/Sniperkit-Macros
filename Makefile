
### BASE_IMAGE #################################################################

# TODO: Use real values
BASE_IMAGE_NAME		?= $(DOCKER_PROJECT)/sniperkit-alpine
BASE_IMAGE_TAG		?= latest

### DOCKER_IMAGE ###############################################################

SAMPLE_PROJECT_VERSION	?= 0.0.1

# TODO: Use real velues
DOCKER_PROJECT		?= roscopecoltran
DOCKER_PROJECT_DESC	?= Sniperkit Docker image
DOCKER_PROJECT_URL	?= https://github.com/sicz/Mk/tree/master/docker-sample-project

DOCKER_NAME			?= sniperkit-project
DOCKER_IMAGE_TAG	?= $(SAMPLE_PROJECT_VERSION)
DOCKER_IMAGE_TAGS	?= latest
DOCKER_IMAGE_DEPENDENCIES += $(SNIPERKIT_CA_IMAGE)

### DOCKER_VERSIONS ############################################################

DOCKER_VERSIONS		?= latest devel

### BUILD ######################################################################

# Docker image build variables
BUILD_VARS		+= SAMPLE_PROJECT_VERSION

# Allows a change of the build/restore targets to the docker-tag if
# the development version is the same as the latest version
DOCKER_CI_TARGET	?= all
DOCKER_BUILD_TARGET	?= docker-build
DOCKER_REBUILD_TARGET	?= docker-rebuild

### DOCKER_EXECUTOR ############################################################

# Use the Docker Compose executor
DOCKER_EXECUTOR		?= compose

# Variables used in the Docker Compose file
COMPOSE_VARS		+= SERVER_CRT_HOST \
			   SNIPERKIT_CA_IMAGE

# Certificate subject aletrnative names
# TODO: Use real values
SERVER_CRT_HOST		+= sniperkit.local

### SNIPERKIT_CA ##################################################################

# Simple CA image
SNIPERKIT_CA_IMAGE_NAME	?= roscopecoltran/sniperkit-ca
SNIPERKIT_CA_IMAGE_TAG		?= latest
SNIPERKIT_CA_IMAGE			?= $(SNIPERKIT_CA_IMAGE_NAME):$(SNIPERKIT_CA_IMAGE_TAG)

# Simple CA service name in the Docker Compose file
SNIPERKIT_CA_SERVICE_NAME	?= $(shell echo $(SNIPERKIT_CA_IMAGE_NAME) | sed -E -e "s|^.*/||" -e "s/[^[:alnum:]_]+/_/g")

# Simple CA container name
ifeq ($(DOCKER_EXECUTOR),container)
SNIPERKIT_CA_CONTAINER_NAME ?= $(DOCKER_EXECUTOR_ID)_$(SNIPERKIT_CA_SERVICE_NAME)
else ifeq ($(DOCKER_EXECUTOR),compose)
SNIPERKIT_CA_CONTAINER_NAME ?= $(DOCKER_EXECUTOR_ID)_$(SNIPERKIT_CA_SERVICE_NAME)_1
else ifeq ($(DOCKER_EXECUTOR),stack)
# TODO: Docker Swarm Stack executor
SNIPERKIT_CA_CONTAINER_NAME ?= $(DOCKER_EXECUTOR_ID)_$(SNIPERKIT_CA_SERVICE_NAME)_1
else
$(error Unknown Docker executor "$(DOCKER_EXECUTOR)")
endif

### MAKE_VARS ##################################################################

# Display the make variables
MAKE_VARS		?= GITHUB_MAKE_VARS \
			   BASE_IMAGE_MAKE_VARS \
			   DOCKER_IMAGE_MAKE_VARS \
			   BUILD_MAKE_VARS \
			   BUILD_TARGETS_MAKE_VARS \
			   EXECUTOR_MAKE_VARS \
			   CONFIG_MAKE_VARS \
			   SHELL_MAKE_VARS \
			   DOCKER_REGISTRY_MAKE_VARS \
			   DOCKER_VERSION_MAKE_VARS


define BUILD_TARGETS_MAKE_VARS
SAMPLE_PROJECT_VERSION:	$(SAMPLE_PROJECT_VERSION)

DOCKER_CI_TARGET:	$(DOCKER_CI_TARGET)
DOCKER_BUILD_TARGET:	$(DOCKER_BUILD_TARGET)
DOCKER_REBUILD_TARGET:	$(DOCKER_REBUILD_TARGET)
endef
export BUILD_TARGETS_MAKE_VARS

define CONFIG_MAKE_VARS
SNIPERKIT_CA_IMAGE_NAME:	$(SNIPERKIT_CA_IMAGE_NAME)
SNIPERKIT_CA_IMAGE_TAG:	$(SNIPERKIT_CA_IMAGE_TAG)
SNIPERKIT_CA_IMAGE:	$(SNIPERKIT_CA_IMAGE)

SERVER_CRT_HOST:	$(SERVER_CRT_HOST)
endef
export CONFIG_MAKE_VARS

### DOCKER_VERSION_TARGETS #####################################################


DOCKER_ALL_VERSIONS_TARGETS ?= build rebuild ci clean

### MAKE_TARGETS ###############################################################

# Build a new image and run the tests
.PHONY: all
all: build clean start wait logs test

# Build a new image and run the tests
.PHONY: ci
ci: $(DOCKER_CI_TARGET)
	@$(MAKE) clean

### BUILD_TARGETS ##############################################################

# Build a new image with using the Docker layer caching
.PHONY: build
build: $(DOCKER_BUILD_TARGET)
	@true

# Build a new image without using the Docker layer caching
.PHONY: rebuild
rebuild: $(DOCKER_REBUILD_TARGET)
	@true

### EXECUTOR_TARGETS ###########################################################

# Display the configuration file
.PHONY: config-file
config-file: display-config-file

# Display the make variables
.PHONY: makevars vars
makevars vars: display-makevars

# Remove the containers and then run them fresh
.PHONY: run up
run up: docker-up

# Create the containers
.PHONY: create
create: docker-create

# Start the containers
.PHONY: start
start: create docker-start

# Wait for the start of the containers
.PHONY: wait
wait: start docker-wait

# Display running containers
.PHONY: ps
ps: docker-ps

# Display the container logs
.PHONY: logs
logs: docker-logs

# Follow the container logs
.PHONY: logs-tail tail
logs-tail tail: docker-logs-tail

# Run shell in the container
.PHONY: shell sh
shell sh: start docker-shell

# Run the tests
.PHONY: test
test: start docker-test

# Run the shell in the test container
.PHONY: test-shell tsh
test-shell tsh:
	@$(MAKE) test TEST_CMD=/bin/bash

# Stop the containers
.PHONY: stop
stop: docker-stop

# Restart the containers
.PHONY: restart
restart: stop start

# Remove the containers
.PHONY: down rm
down rm: docker-rm

# Remove all containers and work files
.PHONY: clean
clean: docker-clean

### MK_DOCKER_IMAGE ############################################################

PROJECT_DIR			?= $(CURDIR)
SNIPERKIT_MK_DIR	?= $(PROJECT_DIR)/../Mk
include $(SNIPERKIT_MK_DIR)/helpers/makefile/docker.image.mk

### GOLANG #####################################################################

fmt:
	go fmt ./...

install-deps:
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/elazarl/go-bindata-assetfs/...

### GOLANG_GLIDE ###############################################################

glide-get:
	go get -v github.com/Masterminds/glide

glide-install:
	glide install --force

### GOLANG_GOM #################################################################

gom-get:
	go get -v github.com/mattn/gom

gom-install:
	glide install --force

### GOLANG_GOX #################################################################

gox-local: glide-get
	gox -verbose -os="$(LOCAL_MACHINE_OS)" -arch="$(LOCAL_MACHINE_ARCH)" -output="/usr/local/sbin/{{.Dir}}" $(glide novendor)

gox-cross: glide-get gox-darwin gox-linux gox-windows
	gox -verbose -os="linux darwin windows" -arch="amd64" -output="/shared/dist/{{.Dir}}/{{.Dir}}_{{.OS}}_{{.ARCH}}" $(glide novendor)

gox-darwin:

gox-linux:

gox-windows:

### GOLANG_FIX #################################################################

logrus-fix:
	rm -fr vendor/github.com/Sirupsen
	find vendor -type f -exec sed -i 's/Sirupsen/sirupsen/g' {} +

### GOLANG_GENERATE ############################################################

generate: clean generate-models

generate-proto:
	protoc --gogofaster_out=. -Iproto -I$(GOPATH)/src proto/caffe2.proto

generate-models:
	go-bindata -nomemcopy -prefix builtin_models/ -pkg caffe2 -o builtin_models_static.go -ignore=.DS_Store  -ignore=README.md builtin_models/...

### GOLANG_CLEAN ###############################################################

clean-models:
	rm -fr builtin_models_static.go

clean-proto:
	rm -fr *pb.go

clean: clean-models

travis: install-deps glide-install logrus-fix generate
	echo "building..."
	go build

################################################################################