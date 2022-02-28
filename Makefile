GO_BIN := go

.PHONY: fetch-deps
fetch-deps: dev-fetch-deps test-fetch-deps

.PHONY: dev-fetch-deps
dev-fetch-deps: test-fetch-deps
	@$(GO_BIN) install github.com/golang/mock/mockgen@v1.6.0

.PHONY: test-fetch-deps
test-fetch-deps:
	@$(GO_BIN) install honnef.co/go/tools/cmd/staticcheck@2021.1.2
	@$(GO_BIN) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.44.2

.PHONY: test-go-version
test-go-version:
	@echo $@
	@./_scripts/assert-go-version.sh

.PHONY: test
test: test-unit test-format test-lint test-security

.PHONY: test-lint
test-lint:
	@echo $@
	@golangci-lint run

.PHONY: test-generate
test-generate:
	@echo $@
	@./_scripts/assert-generated-files-updated.sh

.PHONY: test-format
test-format:
	@echo $@
	@data=$$(gofmt -l .);\
		 if [ -n "$${data}" ]; then \
			>&2 echo "format is broken:"; \
			>&2 echo "$${data}"; \
			exit 1; \
		 fi

.PHONY: test-security
test-security:
	@echo $@
	@staticcheck ./...

.PHONY: test-unit
test-unit:
	@echo $@
	@$(GO_BIN) test ./...
