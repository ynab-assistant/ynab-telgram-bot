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

lint:
	golangci-lint run --config .golangci.yml ./...

generate:
	go generate ./...


bot:
	docker build \
		-f deployment/docker/dockerfile.telegram-bot \
		-t telegram-bot-amd64:1.0 \
		--build-arg VCS_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +”%Y-%m-%dT%H:%M:%SZ”` \
		.