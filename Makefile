swagger:
	swag init -o ./api -g cmd/api/main.go

lint:
	golangci-lint run

test:
	go test -race ./...

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	rm coverage.out

build:
	go build -o ./.bin/minqyy-api ./cmd/api/main.go

run: swagger build
	./.bin/minqyy-api