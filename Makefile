.PHONY: dev build run clean tidy test test-verbose test-coverage test-watch

# Development with hot-reload
dev:
	air

# Build the application
build:
	go build -o ./bin/app ./cmd/main.go

# Run without hot-reload
run:
	go run ./cmd/main.go

# Clean build artifacts
clean:
	rm -rf ./tmp ./bin ./coverage

# Clean up dependencies
tidy:
	go mod tidy

# Run all tests
test:
	go test ./internal/...

# Run tests with verbose output
test-verbose:
	go test -v ./internal/...

# Run tests with coverage report
test-coverage:
	mkdir -p coverage
	go test -coverprofile=coverage/coverage.out ./internal/...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	open coverage/coverage.html

# Run only specific tests matching a pattern
test-filter:
	@read -p "Enter test pattern: " pattern; \
	go test -v ./internal/... -run="$$pattern"

# Run tests in watch mode (requires fswatch: brew install fswatch)
test-watch:
	fswatch -o . | xargs -n1 -I{} make test
