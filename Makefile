generate:
	@echo "Starting..."

	@go run ./cmd/pb-gen/main.go 

build:
	@echo "Building..."

	@go build -o main ./cmd/pb-gen/main.go

start:
	@echo "Starting..."

	@./main
