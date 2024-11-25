build:
	go build -o sr

run1: build
	./sr --port=8090

run2: build
	./sr --port=8091

sample-curl:
	curl -X POST -d @example.json http://localhost:8091/request

docker:
	docker compose up