all: format vet deps test build

format:
	@echo "formatting files..."
	@GO111MODULE=off go get golang.org/x/tools/cmd/goimports
	@goimports -w -l .
	@gofmt -s -w -l .

vet:
	@echo "vetting..."
	@go vet -mod=vendor ./...

deps:
	@echo "installing dependencies..."
	@go get ./...
	@go mod tidy
	@go mod vendor

test:
	@echo "running test coverage..."
	@mkdir -p test-artifacts/coverage
	@go test -mod=vendor ./... -v -coverprofile test-artifacts/cover.out
	@go tool cover -func test-artifacts/cover.out

build: deps
	@echo "building..."
	@go build -mod=vendor ./...
