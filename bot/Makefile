.PHONY:

build:
	go build -o ./bin/bot cmd/bot/main.go

 run: build
	./bin/bot

tidy:
	go mod tidy

test:
	go test --short -coverprofile=cover.out ./...
	make test.coverage

test.coverage:
	go tool cover -func=cover.out