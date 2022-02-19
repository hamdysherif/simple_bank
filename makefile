postgres:
	docker run --name postgres14 -p 5431:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine
createdb:
	docker exec postgres14 createdb -Oroot -Uroot simple_bank
dropdb:
	docker exec postgres14 dropdb simple_bank
migrate-up:
	migrate -path=db/migrate -database="postgresql://root:secret@localhost:5431/simple_bank?sslmode=disable"  -verbose up
migrate-down:
	migrate -path=db/migrate -database="postgresql://root:secret@localhost:5431/simple_bank?sslmode=disable"  -verbose down
sqlc:
	sqlc generate

test:
	go test ./... --cover
