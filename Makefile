test:
	- go test -race ./...

test-cover:
	- go test -race -coverprofile cover.out ./... && go tool cover -html=cover.out -o cover.html && open ./cover.html

clear-cover:
	- rm -rf cover.html cover.out

start-local-dependencies:
	- cd eng && docker-compose up -d

stop-local-dependencies:
	- cd eng && docker-compose down

clear-data-local-dependencies:
	- cd eng && docker-compose down -v

run-web-api:
	- go run cmd/api/api.go

run-worker:
	- go run cmd/worker/worker.go

open-jaeger:
	- open http://localhost:16686/

open-grafana:
	- open http://localhost:3000/



