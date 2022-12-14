ifndef ${TAG}
	TAG=$$(git rev-parse --short HEAD)
endif

packagesToTest=$$(go list ./... | grep -v generated)

setup:
	go install github.com/go-swagger/go-swagger/cmd/swagger@v0.30.0
	go install entgo.io/ent/cmd/ent@v0.11.2
	go install github.com/vektra/mockery/v2@v2.15.0

setup_alpine:
	apk add --update --no-cache git build-base && rm -rf /var/cache/apk/*

run:
	go run ./cmd/swagger/

clean/mocks:
	find ./internal/generated/mocks/* -exec rm -rf {} \; || true

generate/mocks: clean/mocks
	mockery --all --case snake --dir ./pkg/domain --output ./internal/generated/mocks

clean/swagger:
	rm -rf ./internal/generated/swagger

generate/swagger: clean/swagger
	swagger generate server -f ./swagger.yaml -s ./internal/generated/swagger/restapi -m ./internal/generated/swagger/models --exclude-main
	swagger generate client -c ./internal/generated/swagger/client -f ./swagger.yaml -m ./internal/generated/swagger/models

clean/ent:
	find ./internal/generated/ent/* ! -name "generate.go" -exec rm -rf {} \; || true

generate/ent: clean/ent
	go run -mod=mod entgo.io/ent/cmd/ent generate --target ./internal/generated/ent ./internal/ent/schema

generate: generate/swagger generate/ent generate/mocks

clean: clean/swagger clean/ent
	rm csr coverage.out report.txt

build:
	CGO_ENABLED=0 go build -o csr -buildvcs=false ./cmd/swagger/...

lint:
	golangci-lint run --out-format tab | tee ./report.txt

test:
	go test $(go list ./... | grep -v generated) -race -coverprofile=coverage.out -short

coverage:
	go tool cover -func=coverage.out

coverage_total:
	go tool cover -func=coverage.out | tail -n1 | awk '{print $3}' | grep -Eo '\d+(.\d+)?'

int-test:
	DOCKER_BUILDKIT=1  docker build -f ./int-test-infra/Dockerfile.int-test --network host --no-cache -t test_go_run:int-test --target run . && \
	docker-compose -f ./int-test-infra/docker-compose.int-test.yml up -d
	go test -v -timeout 10m ./... -run Integration
	docker-compose -f ./int-test-infra/docker-compose.int-test.yml down

deploy_ssh:
	ssh -o "StrictHostKeyChecking=no" -i ~/.ssh/ssh_deploy -p"${DEPLOY_SSH_PORT}" "${DEPLOY_SSH_USER}@${DEPLOY_SSH_HOST}" 'mkdir -p /var/www/csr/${ENV}/'
	scp -o "StrictHostKeyChecking=no" -i ~/.ssh/ssh_deploy -P"${DEPLOY_SSH_PORT}" -r ./csr "${DEPLOY_SSH_USER}@${DEPLOY_SSH_HOST}:~/tmp_csr"
	scp -o "StrictHostKeyChecking=no" -i ~/.ssh/ssh_deploy -P"${DEPLOY_SSH_PORT}" -r ./config.json "${DEPLOY_SSH_USER}@${DEPLOY_SSH_HOST}:/var/www/csr/${ENV}/"
	ssh -o "StrictHostKeyChecking=no" -i ~/.ssh/ssh_deploy -p"${DEPLOY_SSH_PORT}" "${DEPLOY_SSH_USER}@${DEPLOY_SSH_HOST}" \
	"sudo systemctl daemon-reload && sudo service ${ENV}.csr stop && cp ~/tmp_csr /var/www/csr/${ENV}/server && sudo service ${ENV}.csr start"

