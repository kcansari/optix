.DEFAULT_GOAL := execute
.PHONY: fmt vet build execute clean

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

clean: execute
	@echo "🧹 Cleaning build artifacts..."
	go clean
	@echo "✅ Clean completed"
