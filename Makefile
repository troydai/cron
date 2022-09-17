.PHONY: test
test:
	@go test . -timeout 10s

.PHONY: watch-test
watch-test:
	@watchman-make -p '**/*.go' 'Makefile' -t test


.PHONY: lint
lint:
	@golangci-lint run
