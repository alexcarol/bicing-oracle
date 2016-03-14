.PHONY: fmt lint check test vet check-app docker-up docker-restart

fmt:
	go fmt ./...

lint:
	golint ./...

vet:
	go vet ./...

check: fmt lint vet

test:
	go test -v ./...

check-app:
	wget -p http://localhost/checkup

docker-up:
	docker-compose stop
	docker-compose up -d

docker-restart:
	docker-compose restart
