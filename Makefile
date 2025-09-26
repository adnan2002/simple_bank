postgres:
	sudo docker run --name postgres15 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:15.14-bookworm

createdb:
	sudo docker exec -it postgres15 createdb --username=root --owner=root simple_bank

dropdb:
	sudo docker exec -it postgres15 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

startpostgrescontainer:
	sudo docker start postgres15

server:
	go run .

mock:
	mockgen -package mock -destination db/mock/store.go example.com/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test startpostgrescontainer run mock
