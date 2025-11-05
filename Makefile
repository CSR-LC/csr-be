ifndef ${TAG}
	TAG=$$(git rev-parse --short HEAD)
endif

# Makefile Helpers
GREEN  := \033[0;32m
CYAN   := \033[0;36m
BOLD   := \033[1m
NC     := \033[0m

# Cross-platform open command
OPEN_CMD := open
ifeq ($(shell uname), Linux)
	OPEN_CMD := xdg-open
endif

packagesToTest=$$(go list ./... | grep -v generated)

# Setup and Installation
setup:
	go install github.com/go-swagger/go-swagger/cmd/swagger@v0.30.4
	go install entgo.io/ent/cmd/ent@v0.14.4
	go install github.com/vektra/mockery/v3@v3.2.4

setup_alpine:
	apk add --update --no-cache git build-base && rm -rf /var/cache/apk/*

# Development
run:
	go run ./cmd/swagger/

# Code Generation
clean/mocks:
	find ./internal/generated/mocks/* -exec rm -rf {} \; || true

generate/mocks: clean/mocks
	mockery

clean/swagger:
	cd ./internal/generated/swagger && rm -rfv client models || true
	rm -vf ./internal/generated/swagger/restapi/*.go
	cd ./internal/generated/swagger/restapi/operations && find . ! -name 'gethandlers.go' -type f -exec rm -fv {} +

generate/swagger: clean/swagger
	swagger generate server -P models.Principal -f ./swagger.yaml -s ./internal/generated/swagger/restapi -m ./internal/generated/swagger/models --exclude-main
	swagger generate client -c ./internal/generated/swagger/client -f ./swagger.yaml -m ./internal/generated/swagger/models

clean/ent:
	find ./internal/generated/ent/* ! -name "generate.go" -exec rm -rf {} \; || true

generate/ent: clean/ent
	go run -mod=mod entgo.io/ent/cmd/ent generate --target ./internal/generated/ent ./internal/ent/schema

generate: generate/swagger generate/ent generate/mocks

clean: clean/swagger clean/ent
	rm -f csr coverage.out report.txt

# Build
build:
	CGO_ENABLED=0 go build -o csr ./cmd/swagger/...

# Testing
lint:
	golangci-lint run --out-format tab | tee ./report.txt

test:
	go test ${packagesToTest} -race -coverprofile=coverage.out -short

coverage:
	go tool cover -func=coverage.out

coverage_total:
	go tool cover -func=coverage.out | tail -n1 | awk '{print $3}' | grep -Eo '[0-9]+(\.[0-9]+)?'

# Integration Tests
int-test:
	DOCKER_BUILDKIT=1 docker build -f ./int-test-infra/Dockerfile.int-test --network host --no-cache -t csr:int-test --target run .
	$(MAKE) int-infra-up
	$(MAKE) int-test-without-infra
	$(MAKE) int-infra-down

int-test-without-infra:
	go test -v -p 1 -timeout 10m ./... -run Integration

build-int-image:
	docker build -t csr:int-test -f ./int-test-infra/Dockerfile.int-test .

int-infra-up:
	docker-compose -f ./int-test-infra/docker-compose.int-test.yml up -d --wait

int-infra-down:
	docker-compose -f ./int-test-infra/docker-compose.int-test.yml down

# Database
db:
	docker-compose -f ./docker-compose.yml up -d postgres

schema:
	@echo "==> Generating database schema diagram..."
	@docker run --rm \
		--mount type=bind,source="$$(pwd)",target=/home/schcrwlr \
		--network host \
		schemacrawler/schemacrawler \
		/opt/schemacrawler/bin/schemacrawler.sh \
		--server=postgresql \
		--host=localhost \
		--port=5432 \
		--database=csr \
		--user=csr \
		--password= \
		--schemas=public \
		--info-level=standard \
		--command=schema \
		--output-format=svg \
		--output-file=schema.svg
	@echo ""
	@echo "${GREEN}✔ Diagram generated successfully!${NC}"
	@echo "  ${BOLD}File:${NC}      schema.svg"
	@echo "  ${BOLD}To View:${NC}   Run ${CYAN}${OPEN_CMD} schema.svg${NC}"

# Deployment
deploy_ssh:
	ssh -o "StrictHostKeyChecking=no" -i ~/.ssh/ssh_deploy -p"${deploy_ssh_port}" "${deploy_ssh_user}@${deploy_ssh_host}" 'mkdir -p /var/www/csr/${env}/'
	scp -o "StrictHostKeyChecking=no" -i ~/.ssh/ssh_deploy -P"${deploy_ssh_port}" -r ./csr "${deploy_ssh_user}@${deploy_ssh_host}:~/tmp_csr"
	scp -o "StrictHostKeyChecking=no" -i ~/.ssh/ssh_deploy -P"${deploy_ssh_port}" -r ./config.json "${deploy_ssh_user}@${deploy_ssh_host}:/var/www/csr/${env}/"
	ssh -o "StrictHostKeyChecking=no" -i ~/.ssh/ssh_deploy -p"${deploy_ssh_port}" "${deploy_ssh_user}@${deploy_ssh_host}" \
	"sudo systemctl daemon-reload && sudo service ${env}.csr stop && cp ~/tmp_csr /var/www/csr/${env}/server && sudo service ${env}.csr start"

# Docker Compose Commands
build_project:
	docker-compose build csr

rebuild_project: stop_project
	docker-compose build --no-cache csr
	docker-compose up -d

start_project:
	docker-compose up -d

stop_project:
	docker-compose down

restart_project: stop_project start_project

# Helper targets
.PHONY: help
help:
	@echo "CSR Backend - Available Make Targets"
	@echo ""
	@echo "Setup:"
	@echo "  make setup              - Install required Go tools"
	@echo "  make setup_alpine       - Install dependencies for Alpine Linux"
	@echo ""
	@echo "Development:"
	@echo "  make run                - Run the application locally"
	@echo "  make db                 - Start PostgreSQL database"
	@echo "  make generate           - Generate all code (swagger, ent, mocks)"
	@echo ""
	@echo "Build:"
	@echo "  make build              - Build the application binary"
	@echo ""
	@echo "Testing:"
	@echo "  make test               - Run unit tests"
	@echo "  make lint               - Run linter"
	@echo "  make coverage           - Show test coverage"
	@echo "  make int-test           - Run integration tests"
	@echo ""
	@echo "Docker Compose:"
	@echo "  make build_project      - Build Docker containers"
	@echo "  make rebuild_project    - Rebuild containers from scratch"
	@echo "  make start_project      - Start all services"
	@echo "  make stop_project       - Stop all services"
	@echo "  make restart_project    - Restart all services"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean              - Clean generated files"
