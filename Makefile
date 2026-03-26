all: build

GO_BUILD_FLAGS?=

current_mkfile:=$(abspath $(lastword $(MAKEFILE_LIST)))
current_dir:=$(patsubst %/,%,$(dir $(current_mkfile)))

.PHONY: docs
docs:
	go tool swag init

.PHONY: build
build:
	go build $(GO_BUILD_FLAGS) -a -o $(current_dir)/out/server $(current_dir)/app
