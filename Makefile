build:
	@go build -o bin/splitz-backend cmd/main.go

test:
	@go test -v ./...

run:
	@go run cmd/main.go