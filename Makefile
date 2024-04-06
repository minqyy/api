swagger:
	swag init -g cmd/api/main.go

lint:
	golangci-lint run -E bodyclose -E contextcheck -E dupl -E goconst

test:
	go test -race ./...

build:
	go build -o ./.bin/minqyy-api ./cmd/api/main.go

run: swagger build
	./.bin/minqyy-api