GRPC-DIR=./sipub-tech/grpc

PROTO-FILES=movies.proto

NEED-TEST-MODULES=./sipub-tech/api ./sipub-tech/messaging ./sipub-tech/movies

SOURCE-DIR=./sipub-tech/
MODULES-TO-BUILD=api movies
IMAGE_PREFIX=edmilsonrodrigues/sipub-tech

DEPLOY-DIR=./deploy

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

.PHONY: deploy-microk8s

.PHONY: deploy-cluster-microk8s

.PHONY: deploy-cluster-lxd


