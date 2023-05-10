dev: start_postgres sleep create_db migrate_up

dev_destroy: stop_postgres rm_postgres

server:
	go run main.go

create_bank_network:
	docker network create bank-network

rm_bank_network:
	docker network rm bank-network

start_postgres:
	docker run --name postgres --network bank-network -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d -p 5432:5432 postgres:15.2-alpine3.17

stop_postgres:
	docker stop postgres

rm_postgres:
	docker rm postgres

build_bank:
	docker build -t bank:latest .

start_bank:
	docker run --name bank -it -p 8080:8080 -e GIN_MODE=release --network bank-network -e DB_SOURCE="postgresql://root:secret@postgres:5432/bank?sslmode=disable" bank:latest

stop_bank:
	docker stop bank

rm_bank:
	docker rm bank

create_db:
	docker exec -it postgres createdb --username=root --owner=root bank
	
drop_db:
	docker exec -it postgres dropdb bank
	
migrate_up:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose up

migrate_up_one:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose up 1

migrate_down:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose down

migrate_down_one:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose down 1
	
sqlc_generate:
	sqlc generate

mockgen:
	mockgen -package mockdb -destination db/mock/store.go github.com/CrunchyBlue/Golang-Bank/sqlc Store
	
test:
	go test -v -cover ./...
	
sleep:
	sleep 3
	
.PHONY: dev dev_destroy server create_bank_network start_postgres stop_postgres rm_postgres build_bank start_bank create_db drop_db migrate_up migrate_down migrate_down_one sqlc_generate mockgen test sleep