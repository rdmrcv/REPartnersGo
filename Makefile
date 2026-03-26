all: build

GO_BUILD_FLAGS?=

current_mkfile:=$(abspath $(lastword $(MAKEFILE_LIST)))
current_dir:=$(patsubst %/,%,$(dir $(current_mkfile)))

current_dir_name:=$(notdir $(current_dir))

.PHONY: build
build:
	go build $(GO_BUILD_FLAGS) -mod=vendor -a -o $(current_dir)/out/$(current_dir_name) $(current_dir)
