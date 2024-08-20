lint:
	gofumpt -w .
	gci write . --skip-generated -s standard -s default
	golangci-lint run -v --fix --timeout=5m ./...
	go mod tidy