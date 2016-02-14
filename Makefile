.PHONY: fmt lint check test vet

fmt:
	go fmt ./...

lint:
	golint ./...

vet:
	go vet ./...

check: fmt lint vet

test:
	go test -v ./...
