# For tests you will need docker API on port 2375 
# with disabled tls
cover:
	go test -coverprofile cover.out ./...
	go tool cover -html cover.out

test:
	go test ./... -v

build:
	go build -o bin/ cmd/app/main.go

run:
	go run cmd/app/main.go 1> logs.log

swagger:
	swag fmt
	swag init --parseDependency -g internal/controller/http/v1/rest/word.go

migrateup:
	migrate -database postgresql://flash_cards:12345@localhost:5432/word_api?sslmode=disable -path ./internal/repository/postgresql/migrations/ up 

migratedown:
	migrate -database postgresql://flash_cards:12345@localhost:5432/word_api?sslmode=disable -path ./internal/repository/postgresql/migrations/ down 

