.PHONY: test
test:
	go test ./... -v -coverprofile=coverprofile

.PHONY: test-cover
test-cover:
	go test ./... -coverprofile=coverprofile

.PHONY: test-all
test-all:
	go test ./... -v -tags=integration -coverprofile=coverprofile

.PHONY: test-all-cover
test-all-cover:
	go test ./... -coverprofile=coverprofile

