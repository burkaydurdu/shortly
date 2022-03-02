.PHONY: lint
lint:
	golangci-lint run -c .golangci.yml -v --fix

.PHONY: unit-test
unit-test:
	go test -v ./... -coverprofile=unit_coverage.out -short -tags=unit