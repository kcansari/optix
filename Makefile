.DEFAULT_GOAL := execute
.PHONY: fmt vet build execute clean debug test-coverage test-verbose debug-build run-debug setup-test-data

## Formats all Go files in the current directory and subdirectories
fmt:
	@echo "🔧 Formatting Go code..."
	go fmt ./...
	@echo "✅ Code formatting completed"

## checks the code for potential issues
vet: fmt
	@echo "🔍 Running Go vet analysis..."
	go vet ./...
	@echo "✅ Code analysis completed"

## compiles the Go program
build: vet
	@echo "🏗️  Building Go program..."
	go build -o optix
	@echo "✅ Build completed successfully"

## executes the Go program
execute: build
	@echo "🚀 Executing the program..."
	./optix

clean:
	@echo "🧹 Cleaning build artifacts..."
	go clean
	rm -f optix
	rm -f debug
	@echo "✅ Clean completed"

## Development and Debugging Targets

## builds the binary with debug information
debug-build:
	@echo "🔧 Building debug version..."
	go build -gcflags="-N -l" -o debug
	@echo "✅ Debug build completed"

## runs tests with verbose output
test-verbose:
	@echo "🧪 Running tests with verbose output..."
	go test -v ./...
	@echo "✅ Verbose tests completed"

## runs tests with coverage report
test-coverage:
	@echo "📊 Running tests with coverage..."
	go test -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

