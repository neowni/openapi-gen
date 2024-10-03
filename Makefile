
install:
	go install ./cmd/n-cl

lint:
	golangci-lint run -c .golangci.yaml ./...
