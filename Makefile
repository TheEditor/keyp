.PHONY: build test clean test-coverage install

# Build binary
build:
	go build -o keyp ./cmd/keyp

# Run all tests
test:
	go test ./... -v

# Run tests with coverage
test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

# Clean build artifacts
clean:
	rm -f keyp keyp.exe coverage.out

# Install locally
install: build
	cp keyp /usr/local/bin/keyp
