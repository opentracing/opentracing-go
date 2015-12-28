

.DEFAULT_GOAL := test

.PHONY: test
test:
	go test -v -cover ./...

.PHONY: lint
lint:
	golint ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: example
example:
	go build -o build/dapperish-example ./examples/dapperish.go
	./build/dapperish-example
