.PHONY: build build-all clean test run

# Build for current platform
build:
	go build -o frontforge

# Build for all platforms
build-all: clean
	@echo "Building for multiple platforms..."
	GOOS=darwin GOARCH=amd64 go build -o bin/frontforge-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o bin/frontforge-darwin-arm64
	GOOS=linux GOARCH=amd64 go build -o bin/frontforge-linux-amd64
	GOOS=windows GOARCH=amd64 go build -o bin/frontforge-windows-amd64.exe
	@echo "Build complete! Binaries are in ./bin/"

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f frontforge frontforge.exe

# Run the CLI
run:
	go run main.go

# Run tests (when tests are added)
test:
	go test -v ./...
