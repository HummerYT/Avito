.PHONY: up

up:
	docker build -t avito .
	docker compose up -d app

down:
	docker compose down

unit:
	go test ./...

cover:
	go test -coverprofile="coverage.out" ./...
	go tool cover -func="coverage.out"