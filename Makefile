.PHONY: all install-mockery mocks test lint fmt coverage check

all: check

install-mockery:
	go install github.com/vektra/mockery/v3@v3.5.1

mocks:
	mockery

test:
	go test -v ./... -coverprofile=coverage.out -coverpkg=./internal/...

lint:
	go vet ./...

fmt:
	go fmt ./...

coverage:
	go tool cover -func=coverage.out

check: fmt lint test coverage
