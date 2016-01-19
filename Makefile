docker-build:
	docker build -t bicing_stats .

linux-build-query:
	env GOOS=linux go build -o querybicing tools/querybicing.go && docker build -t alexcarol/querybicing -f Dockerfile.querybicing .

docker-up: docker-build
	docker run -it --rm --name bicing_app bicing_stats

test:
	go test -v ./...
