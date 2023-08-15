.PHONY: all tidy generate lint vet test coverage pushdoc

# Default "make" target to check locally that everything is ok, BEFORE pushing remotely
all: lint vet test build
	@echo "Done with the standard checks"

tidy:
	# Tidy up go modules
	@go mod tidy

# Some packages are excluded from staticcheck due to deprecated warnings: #208.
lint: tidy
	# Static checks
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.49.0
	golangci-lint run

vet: tidy
	# Go vet
	@go vet ./...

test: tidy
	# Test without coverage
	@go test ./...

build: tidy
	@make -C smc/smccli
	@make -C blockchain/chaincli

coverage: tidy
	# Test and generate a coverage output usable by sonarcloud
	@go test -json -covermode=count -coverprofile=profile.cov ./... | tee report.json

