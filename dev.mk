### Manage Dependencies

## Install dependencies
deps.install:
	# install golanglint-ci into ./bin
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.41.1
	# Download packages in go.mod file
	go mod download
.PHONY: deps.install

## Update dependencies
deps.update:
	# Update go dependencies in go.mod and go.sum
	go get -u ./cmd/...
	go get -u ./internal/...
	go mod tidy
.PHONY: deps.update

### Code verification and static analysis

## Run code verification
verify:
	# Lint go files
	./bin/golangci-lint --version
	./bin/golangci-lint --timeout 3m0s run ./...
.PHONY: verify

## Run code verification and autofix issues where possible
verify.fix:
	echo 'TODO: add go code verification autofixing where possible'
.PHONY: verify.fix

### Build

## Build binary
build:
	mkdir -p ./build
	env CGO_ENABLED=0 go build -o build/morningpaper2remarkable ./cmd/morningpaper2remarkable/main.go
.PHONY: build

### Testing

## Run tests
test: test.unit
.PHONY: test

## Run tests and output reports
test.report: test.unit.report
.PHONY: test.report

## Run unit tests
test.unit:
	go test -count=1 -v -p=1  ./internal/... ./cmd/...
.PHONY: test.unit

## Run unit tests and output reports
test.unit.report:
	mkdir -p reports
	go test -json -count=1 -coverprofile=reports/test-unit.out -v -p 5 ./internal/... ./cmd/... > reports/test-unit.json
.PHONY: test.unit.report
