createdb:
	@psql -U okolomichael -d postgres -c 'CREATE DATABASE simple_bank'

dropdb:
	@psql -U okolomichael -d postgres -c 'DROP DATABASE simple_bank'

migrateup:
	@migrate -path db/migration -database "postgresql://okolomichael:Nonso007@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	@migrate -path db/migration -database "postgresql://okolomichael:Nonso007@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	@sqlc generate

test:
	@go test -v -cover ./...

.PHONY: createdb dropdb migrateup migratedown sqlc test