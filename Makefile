PKG_NAME:=github.com/sapcc/lyra-cli
BUILD_IMAGE:=golang:1.16
TARGETS:=linux/amd64 windows/amd64 darwin/amd64 darwin/arm64

LYRA_CLI_BIN_TPL:=lyra_cli_{{.OS}}_{{.Arch}}
ifneq ($(BUILD_VERSION),)
LDFLAGS += -X github.com/sapcc/lyra-cli/version.Version=$(BUILD_VERSION)
LYRA_CLI_BIN_TPL:=lyra_cli_$(BUILD_VERSION)_{{.OS}}_{{.Arch}}
endif

.PHONY: help
help:
	@echo
	@echo "Available targets:"
	@echo "  * test              - run metalint and unit"
	@echo "  * unit              - run unit tests"
	@echo "  * metalint          - linter which runs a number of other linters against your files, and normalises their output to a standard format."
	@echo "  * cross             - cross compile for darwin, windows, linux (requires docker)"

packages = $(PKG_NAME) $(shell go list -f '{{ join .Deps "\n" }}' | grep -v vendor | grep $(PKG_NAME))

.PHONY: test
test: metalint unit

.PHONY: unit
unit:
	go test -v -timeout=120s ./...

.PHONY: metalint
metalint:
	docker run --rm -v $(CURDIR):/app -w /app golangci/golangci-lint:v1.41.1 golangci-lint run -v

.PHONY: cross
cross:
	@# -w omit DWARF symbol table -> smaller
	@# -s stip binary
	docker run \
		--rm \
		-v $(CURDIR):/go/src/$(PKG_NAME) \
		-w /go/src/$(PKG_NAME) \
		$(BUILD_IMAGE) \
		go get github.com/mitchellh/gox && make cross-compile TARGETS="$(TARGETS)" BUILD_VERSION=$(BUILD_VERSION)

.PHONY: cross-compile
cross-compile:
	gox -osarch="$(TARGETS)" -output="bin/$(LYRA_CLI_BIN_TPL)" -ldflags="$(LDFLAGS)" $(PKG_NAME)
