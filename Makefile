
.PHONY: install
install:
	go install ./cmd/n-cl

.PHONY: lint
lint:
	golangci-lint run -c .golangci.yaml ./...

.PHONY: test
test:
	go run test/main.go
