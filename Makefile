.PHONY: fmt lint check test

fmt:
	go fmt ./...

lint:
	golint ./...

check: fmt lint

test:
	go test -v ./...
