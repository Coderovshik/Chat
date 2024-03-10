CONTAINER_NAME=postgres-chat-db

.PHONY: build-server
build-server:
	@go build -o bin/server cmd/server/main.go

.PHONY: run-server
run-server: build-server
	@./bin/server

.PHONY: build-migrator
build-migrator:
	@go build -o bin/migrator cmd/migrator/main.go

.PHONY: run-migrator
run-migrator:
	@./bin/migrator

.PHONY: clean
clean:
	@go rm -rf bin

.PHONY: postgres-init
postgres-init:
	@docker run --name ${CONTAINER_NAME} -p 5433:5432 \
	-e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -d postgres:16-alpine

.PHONY: postgres-rm
postgres-rm:
	@docker stop ${CONTAINER_NAME} && docker rm ${CONTAINER_NAME}

.PHONY: postgres-cli
postgres-cli:
	@docker exec -it ${CONTAINER_NAME} psql -U admin

.PHONY: create-db
create-db:
	@docker exec -it ${CONTAINER_NAME} createdb --username=admin --owner=admin chat

.PHONY: drop-db
drop-db:
	@docker exec -it ${CONTAINER_NAME} dropdb chat