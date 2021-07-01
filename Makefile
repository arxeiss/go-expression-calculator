.PHONY: run test tests fix_lint lint install_tools

run:
	@go run .

test: tests

tests:
	go test -v ./...

lint:
	golangci-lint run ./...

build:
	go build -o calculator

fix_lint:
	golangci-lint run --fix

install_tools:
	@echo Installing tools from tools.go
	go list -f '{{range .Imports}}{{.}} {{end}}' tools.go | xargs go install
