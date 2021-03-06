postgres:
	docker run --name postgres14 -network=simplebank-network -p 5431:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine
createdb: |
	docker exec postgres14 createdb -Oroot -Uroot simple_bank
	docker exec postgres14 createdb -Oroot -Uroot simple_bank_test
dropdb: |
	docker exec postgres14 dropdb simple_bank
	docker exec postgres14 dropdb simple_bank_test
migrate-up: |
	migrate -path=db/migrate -database="postgresql://root:secret@localhost:5431/simple_bank?sslmode=disable"  -verbose up ${v}
	migrate -path=db/migrate -database="postgresql://root:secret@localhost:5431/simple_bank_test?sslmode=disable"  -verbose up ${v}
migrate-down: |
	migrate -path=db/migrate -database="postgresql://root:secret@localhost:5431/simple_bank?sslmode=disable"  -verbose down ${v}
	migrate -path=db/migrate -database="postgresql://root:secret@localhost:5431/simple_bank_test?sslmode=disable"  -verbose down ${v}

#usage make migrate-create name=[migration_name]
migrate-create:
	migrate create -ext sql -dir db/migrate -seq ${name}
sqlc:
	sqlc generate
test:
	go test ./... --cover
server:
	go run main.go
mockgen:
	mockgen -package mockdb -destination db/mock/store.go github.com/hamdysherif/simplebank/db/sqlc Store

.PHONY: postgres createdb dropdb migrate-up migrate-down sqlc test server mockgen
