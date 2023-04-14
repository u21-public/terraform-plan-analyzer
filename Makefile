SRC_ROOT=.

.PHONY: setup
setup:
	@echo "Setting up tools..."
	@test -x ${GOPATH}/bin/golangci-lint || go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.2

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	gofmt -s -w ${SRC_ROOT}

.PHONY: lint
lint: setup
	@echo "Linting code..."
	golangci-lint -v run --timeout 3m $(if $(filter true,$(fix)),--fix,)

.PHONY: tidy
tidy:
	@echo "Fetching dependencies..."
	go mod tidy

.PHONY: test
test: tidy
	@echo "Running tests..."
	go test -v -race -short -cover -coverprofile cover.out ${SRC_ROOT}/... -tags integration
	go tool cover -func cover.out

.PHONY: build
build: tidy
	@echo "Building binary..."
	go build -o ./bin/${BIN_NAME} ${SRC_ROOT}/cmd/main.go
