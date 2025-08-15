GRPC-DIR=./sipub-tech/grpc
PROTO-FILE=movies.proto

.PHONY: update-proto
update-proto:
	protoc -I=${GRPC-DIR} --go_out=${GRPC-DIR} ${GRPC-DIR}/${PROTO-FILE}
