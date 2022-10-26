mod:
	@go mod tidy
	@go mod vendor

lint:
	@go fmt ./...
	@golangci-lint run --verbose --timeout 5m0s --config ./.golangci.yml

up:
	@go run ./cmd/app/... 

