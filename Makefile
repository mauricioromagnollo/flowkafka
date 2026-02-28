# ================================================
# VARIABLES
# ================================================
APP_NAME = kafrest
ENV_FILE = --env-file .env


# ================================================
# COMMANDS
# ================================================
help: ## Print all available commands
	$(info ========================================)
	$(info Available Commands:)
	@grep '^[[:alnum:]_-]*:.* ##' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS=":.* ## "}; {printf "make %-25s %s\n", $$1, $$2};'
	$(info ========================================)
.PHONY: help

up: ## Up all containers in detached mode
	@docker compose up -d
.PHONY: up

down: ## Down all containers, but keep images, networks, and volumes
	@docker compose down
.PHONY: down

stop: ## Stop all containers, but keep images, networks, and volumes
	@docker compose stop
.PHONY: stop

clear: ## Remove all containers, images, networks, and volumes
	@docker compose down --rmi all --volumes --remove-orphans
.PHONY: clear

test: ## Run all tests
	@echo "Running tests..."
	go test -v ./...
.PHONY: test

test-coverage: ## Run test coverage
	go test -v -race -coverprofile=cover.out ./...
	go tool cover -func=cover.out
.PHONY: test-coverage

test-coverage-web: ## Run test coverage and show in browser
	go test -v -race -coverprofile=cover.out ./... && go tool cover -html=cover.out
.PHONY: test-coverage-web

test-race: ## Run data race tests
	CGO_ENABLED=1 go test -race ./...
.PHONY: test-race