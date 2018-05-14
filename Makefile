PKG_NAME:=github.com/sapcc/lyra-cli

.PHONY: help
help:
	@echo
	@echo "Available targets:"
	@echo "  * test              - run metalint and unit"
	@echo "  * unit              - run unit tests"
	@echo "  * metalint          - linter which runs a number of other linters against your files, and normalises their output to a standard format."

packages = $(PKG_NAME) $(shell go list -f '{{ join .Deps "\n" }}' | grep -v vendor | grep $(PKG_NAME))

.PHONY: test
test: metalint unit

.PHONY: unit
unit:
	go test -v -timeout=4s ./...

.PHONY: metalint
metalint:
	gometalinter --vendor --disable-all -E goimports -E megacheck -E ineffassign -E gas --deadline=60s ./...
