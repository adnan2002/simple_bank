include app.env
export

postgres:
	sudo docker run --name postgres15 --network bank-network -e POSTGRES_USER=root -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:15.14-bookworm

createdb:
	sudo docker exec -it postgres15 createdb --username=root --owner=root simple_bank

dropdb:
	sudo docker exec -it postgres15 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)" -verbose down 1

migrateup1:
	migrate -path db/migration -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)" -verbose up 1

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

proto:
	rm -rf pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb  --grpc-gateway_opt paths=source_relative \
    proto/*.proto
evans:
	evans --host localhost --port 9090 --reflection repl

.PHONY: postgres evans createdb dropdb migrateup migratedown sqlc test startpostgrescontainer server mock migratedown1 migrateup1 proto