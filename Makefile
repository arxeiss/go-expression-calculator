.PHONY: run test tests fix_lint lint install_tools

run:
	@go run .

test: tests

tests:
	go test -v -covermode=count -coverpkg github.com/arxeiss/go-expression-calculator/... -coverprofile=coverage.full.out ./...
	# Remove cmd package from coverage, see https://dev.to/arxeiss/false-positive-go-code-coverage-3k7j
	cat coverage.full.out | grep -v "go-expression-calculator/cmd/" > coverage.out
	rm -f coverage.full.out

coverage: tests
	go tool cover -func coverage.out

lint:
	golangci-lint run ./...

build:
	go build -o calculator

fix_lint:
	golangci-lint run --fix

install_tools:
	@echo Installing tools from tools.go
	go list -f '{{range .Imports}}{{.}} {{end}}' tools.go | xargs go install
