# For tests you will need docker API on port 2375 
# with disabled tls
cover:
	go test -coverprofile cover.out ./...
	go tool cover -html cover.out

test:
	go test ./...

build:
	go build -o bin/ cmd/app/main.go

run:
	go run cmd/app/main.go