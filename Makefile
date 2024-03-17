.PHONY: migration-up migration-down lint test

migration-up:
	migrate -path database/migration/ -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" -verbose up

migration-down:
	migrate -path database/migration/ -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" -verbose down

lint:
	golangci-lint run ./...

test:
	go test ./... -short