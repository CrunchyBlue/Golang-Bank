build: start_postgres sleep create_db migrate_up

destroy: stop_postgres rm_postgres

server:
	go run main.go

start_postgres:
	docker run --name postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d -p 5432:5432 postgres:alpine

stop_postgres:
	docker stop postgres

rm_postgres:
	docker rm postgres

create_db:
	docker exec -it postgres createdb --username=root --owner=root bank
	
drop_db:
	docker exec -it postgres dropdb bank
	
migrate_up:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose down
	
sqlc_generate:
	sqlc generate
	
test:
	go test -v -cover ./...
	
sleep:
	sleep 3
	
.PHONY: build destroy server start_postgres stop_postgres rm_postgres create_db drop_db migrate_up migrate_down sqlc_generate test sleep