

.DEFAULT_GOAL := test

.PHONY: test
test:
	go test -v -cover ./...

.PHONY: lint
lint:
	go fmt ./...
	golint ./...
	@# Run again with magic to exit non-zero if golint outputs anything.
	@! (golint ./... | read dummy)
	go vet ./...

