

.DEFAULT_GOAL := test

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golint ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: example
example:
	go build -o build/dapperish-example ./examples/dapperish/.
	./build/dapperish-example
