PKG_NAME:=github.com/sapcc/lyra-cli
BUILD_IMAGE:=hub.***REMOVED***/monsoon/gobuild:1.10
TARGETS:=linux/amd64 windows/amd64

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
	go test -v -timeout=20s ./...

.PHONY: metalint
metalint:
	gometalinter --vendor --disable-all -E goimports -E megacheck -E ineffassign -E gas --deadline=60s ./...

.PHONY: cross
cross:
	@# -w omit DWARF symbol table -> smaller
	@# -s stip binary
	docker run \
		--rm \
		-v $(CURDIR):/go/src/$(PKG_NAME) \
		-w /go/src/$(PKG_NAME) \
		$(BUILD_IMAGE) \
		make cross-compile TARGETS="$(TARGETS)" BUILD_VERSION=$(BUILD_VERSION)

.PHONY: cross-compile
cross-compile:
	gox -osarch="$(TARGETS)" -output="bin/$(LYRA_CLI_BIN_TPL)" -ldflags="$(LDFLAGS)" $(PKG_NAME)
