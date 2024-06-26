.PHONY: all tidy build install lint vet test coverage

# Default "make" target to check locally that everything is ok, BEFORE pushing remotely
all: lint vet test build
	@echo "Done with the standard checks"

tidy:
	# Tidy up go modules
	@go mod tidy

# Some packages are excluded from staticcheck due to deprecated warnings: #208.
lint: tidy
	# Static checks
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint run

vet: tidy
	# Go vet
	@go vet ./...

test: tidy
	# Test without coverage
	@go test ./...

build: tidy
	# Building executable
	@go build ./smc/smccli
	@go build ./blockchain/chaincli
	@go build ./registry/registrycli

install: tidy
	# Building executable
	@go install ./smc/smccli
	@go install ./blockchain/chaincli
	@go install ./registry/registrycli

coverage: tidy
	# Test and generate a coverage output usable by sonarcloud
	@go test -json -covermode=count -coverpkg=./... -coverprofile=profile.cov ./... | tee report.json

