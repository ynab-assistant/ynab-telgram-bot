.PHONY:

build:
	go mod download && CGO_ENABLED=0 go build -o ./bin/ynab-bot cmd/bot/main.go

 run: build
	./bin/ynab-bot

tidy:
	go mod tidy

test:
	go test --short -coverprofile=cover.out ./...
	make test.coverage

test.coverage:
	go tool cover -func=cover.out

test.coverage.html:
	go tool cover -html=cover.out