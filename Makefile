

.DEFAULT_GOAL := test

.PHONY: test
test:
	go test -v -cover ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golint ./...
	@# Run again with magic to exit non-zero if golint outputs anything.
	@! (golint ./... | read dummy)

.PHONY: vet
vet:
	go vet ./...

