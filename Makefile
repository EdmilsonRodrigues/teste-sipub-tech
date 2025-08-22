GRPC-DIR=./sipub-tech/grpc

PROTO-FILES=movies.proto

NEED-TEST-MODULES=./sipub-tech/api ./sipub-tech/messaging ./sipub-tech/movies

SOURCE-DIR=./sipub-tech/
MODULES-TO-BUILD=api movies
IMAGE_PREFIX=edmilsonrodrigues/sipub-tech

DEPLOY-DIR=./deploy

FAKE_AWS_REGION = "us-east-1"
JSON_PATH ?= data/movies.json
DYNAMO_DB_ENDPOINT ?= http://localhost:4566


.PHONY: update-proto
update-proto:
	for file in $(PROTO-FILES); do \
		protoc -I=${GRPC-DIR} --go-grpc_out=${GRPC-DIR} --go_out=${GRPC-DIR} ${GRPC-DIR}/$${file}; \
	done

.PHONY: lint
lint:
	for module in $(NEED-TEST-MODULES); do \
		golangci-lint run $${module}/...; \
	done

.PHONY: test
test:
	for module in $(NEED-TEST-MODULES); do \
		go test $${module}/...; \
	done

.PHONY: build
build:
	for module in $(MODULES-TO-BUILD); do \
		docker image build ${SOURCE-DIR}/$${module} -t ${IMAGE_PREFIX}-$${module}; \
	done

.PHONY: deploy-docker
deploy-docker:
	cd sipub-tech && \
	docker compose up --build

.PHONY: deploy-microk8s
	cd deploy/k8s-local && \
	./setup.bash

.PHONY: deploy-lxd
deploy-lxd:
	cd deploy/k8s-local && \
	./deploy.bash

.PHONY: fill-db
fill-db:
	JSON_PATH=${JSON_PATH} DYNAMO_DB_ENDPOINT=${DYNAMO_DB_ENDPOINT} AWS_REGION=${FAKE_AWS_REGION} go run sipub-tech/movies/tools/importer.go

