ifeq ($(wildcard .env),)
    $(shell cp .env.example .env)
endif

include .env
export

.PHONY: help
help: ## Show this menu
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: clean
clean: ## Clean all temp files
	@rm -rf tmp/

.PHONY: install-migrate
install-migrate: ## install migrate
	@if [ ! -f /usr/local/bin/migrate ]; then \
		wget -O /tmp/migrate.tar.gz https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz; \
		tar -C /tmp -xzvf /tmp/migrate.tar.gz; \
		sudo mv /tmp/migrate /usr/local/bin/migrate; \
	else \
		echo "Great, you already have [migrate] installed"; \
	fi

.PHONY: install-evans
	@if [ ! -f /usr/local/bin/evans ]; then \
		wget -O /tmp/evans.tar.gz https://github.com/ktr0731/evans/releases/download/v0.10.11/evans_linux_amd64.tar.gz; \
		tar -xzvf /tmp/evans.tar.gz; \
		sudo mv evans /usr/local/bin; \
	else
		echo "Great, you already have [evans] installed"; \
	fi

.PHONY: install
install: install-migrate install-evans ## install all requirements

.PHONY: migration_init
migration_init: ## init migrations
	@migrate create -ext=sql -dir=internal/infra/sql/migrations -seq init

.PHONY: migration_up
migration_up: ## run migrations up
	@migrate -path=internal/infra/sql/migrations -database "mysql://root:${DB_PASSWORD}@tcp(localhost:3306)/${DB_NAME}" -verbose up
	@docker exec -it -e MYSQL_PASSWORD=root mysql mysql -uroot -p${DB_PASSWORD} -D ${DB_NAME} -e "show tables;"

.PHONY: migration_down
migration_down: ## run migrations down
	@migrate -path=internal/infra/sql/migrations -database "mysql://root:${DB_PASSWORD}@tcp(localhost:3306)/${DB_NAME}" -verbose down --all
	@docker exec -it -e MYSQL_PASSWORD=root mysql mysql -uroot -p${DB_PASSWORD} -D ${DB_NAME} -e "show tables;"

.PHONY: migration_clean
migration_clean: migration_down migration_up ## run migration down and up to cleanup all data

.PHONY: test
test: ## Run unit-tests
	@go test -v ./...

.PHONY: build
build: ## Build the container image
	@docker build -t github.com/fabiowgermano/clean-architecture:dev -f Dockerfile .

.PHONY: up
up: ## Put the compose containers up
	@docker-compose up -d

.PHONY: down
down: ## Put the compose containers down
	@docker-compose down

.PHONY: serve
serve: ## Run the rest-api, grpc and graphql servers
	@cd cmd/ordersystem && go run main.go wire_gen.go
	@cd -