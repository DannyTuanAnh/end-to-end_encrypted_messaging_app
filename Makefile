include .env
export

CONN_STRING = "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)"

MIGRATION_DIRS = ./internal/db/migrations

# import database
import-db:
	docker exec -i postgres-db psql -U root -d master-golang < ./backupdb-master-master-golang.sql

# export database
export-db:
	docker exec -i postgres-db pg_dump -U root -d master-golang > ./backupdb-master-master-golang.sql

# create new migration file (Example: make migrate-create name=profiles)
migrate-create:
	migrate create -ext sql -dir $(MIGRATION_DIRS) -seq $(name) 

# run all up migrations (Example: make migrate-up)
migrate-up: 
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" up

# rollback the last migration (Example: make migrate-down)
migrate-down: 
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" down 1

# rollback n migration (Example: make migrate-down n=2)
migrate-down-n: 
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" down $(n)

# force migration version (use with caution, Example: make migrate-force version=2) 
migrate-force:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" force $(version)

# drop everything (include schema migration)
migrate-drop:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" drop 

# apply specific migration version (Example: make migrate-goto version=2)
migrate-goto:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" goto $(version)

# Generate Go code from SQL queries using sqlc
sqlc:
	sqlc generate

buf:
	buf dep update

# Makefile for generating gRPC code from .proto files
# Example: make gen-proto path=calculator/calculatorpb/calculator.proto
# path is the path to your .proto file
gen-proto:
	protoc $(path) --go_out=. --go-grpc_out=.

run-api:
	go run cmd/api/main.go

run-auth:
	go run cmd/auth/main.go

run-user:
	go run cmd/user/main.go

.PHONY: import-db export-db migrate-create migrate-up migrate-down migrate-force migrate-drop migrate-goto sqlc gen-proto run-api run-auth run-user buf
