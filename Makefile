.PHONY: vendor

# download vendor
vendor:
	GO111MODULE=on go mod vendor

# run dev chat server
run:
	docker-compose down
	cp ./.env.template ./.env
	docker-compose up --build

# start all unit tests
tests:
	go test ./app/...

# build chat server
build:
	go build -mod=vendor -o bin/chat_server ./cmd/...

# start linter
lint:
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./app/...
