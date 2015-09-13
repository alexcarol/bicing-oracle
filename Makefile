docker-build:
	docker build -t bicing_stats .

docker-up: docker-build
	docker run -it --rm --name bicing_app bicing_stats

test:
	go test -v ./...
