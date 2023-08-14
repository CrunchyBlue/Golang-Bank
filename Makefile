aws_get_secrets:
	aws secretsmanager get-secret-value --secret-id golang-bank --region us-east-1 --profile david --query SecretString --output text | jq -r 'to_entries|map("\(.key)=\(.value)")|.[]' > app.env

aws_sso_login:
	aws sso login --profile david

build_bank:
	docker build -t bank:latest .

create_bank_network:
	docker network create bank-network

create_db:
	docker exec -it postgres createdb --username=root --owner=root bank

dev: start_postgres sleep create_db migrate_up

dev_destroy: stop_postgres rm_postgres
	
drop_db:
	docker exec -it postgres dropdb bank

generate_symmetric_key:
	openssl rand -hex 64 | head -c 32

install_cert_manager:
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml

install_ingress_controller:
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.7.1/deploy/static/provider/aws/deploy.yaml

k9s:
	k9s

migrate_down:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose down

migrate_down_one:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose down 1
	
migrate_up:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose up

migrate_up_one:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose up 1

mockgen:
	mockgen -package mockdb -destination db/mock/store.go github.com/CrunchyBlue/Golang-Bank/sqlc Store

server:
	go run main.go

rm_bank:
	docker rm bank

rm_bank_network:
	docker network rm bank-network

rm_postgres:
	docker rm postgres
	
sleep:
	sleep 3

start_bank:
	docker run --name bank -it -p 8080:8080 -e GIN_MODE=release --network bank-network -e DB_SOURCE="postgresql://root:secret@postgres:5432/bank?sslmode=disable" bank:latest

start_postgres:
	docker run --name postgres --network bank-network -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d -p 5432:5432 postgres:15.2-alpine3.17

stop_bank:
	docker stop bank

stop_postgres:
	docker stop postgres
	
sqlc_generate:
	sqlc generate
	
test:
	go test -v -cover ./...

update_kubeconfig:
	aws eks update-kubeconfig --name golang-bank --region us-east-1
	
.PHONY: aws_get_secrets aws_sso_login build_bank create_bank_network create_db dev dev_destroy drop_db generate_symmetric_key install_cert_manager install_ingress_controller k9s migrate_down migrate_down_one migrate_up migrate_up_one mockgen server rm_bank rm_bank_network rm_postgres sleep start_bank start_postgres stop_bank stop_postgres sqlc_generate test update_kubeconfig
