docker-build:
	docker build -t bicing_stats .

linux-build-query:
	env GOOS=linux GOARCH=amd64 go build -o querybicing_amd64_v3 tools/querybicing.go

docker-up: docker-build
	docker run -it --rm --name bicing_app bicing_stats

test:
	go test -v ./...
