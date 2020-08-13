.PHONY: build

# download vendor
vendor:
	GO111MODULE=on go mod vendor

# run dev chat server
run:
	docker-compose down
	cp ./configs/.env.template ./.env
	docker-compose up --build

# start all unit tests
tests:
	go test ./...

# build chat server
build:
	go build -mod=vendor -o bin/webitel.chat.service.storage ./chat_storage/*.go
	go build -mod=vendor -o bin/webitel.chat.service.api ./chat_api/*.go
	go build -mod=vendor -o bin/webitel.chat.service.telegrambot ./telegram_bot/*.go
	go build -mod=vendor -o bin/webitel.chat.service.flowadapter ./flow_adapter/*.go

# start linter
lint:
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./...

# cd & generate boiler models
generate-boiler:
	sqlboiler --wipe --no-tests -o ./models -c ./configs/sqlboiler.toml psql

proto: 
	protoc --proto_path=. --go_out=. --micro_out=.  chat_storage/proto/storage/storage.proto
	protoc --proto_path=. --go_out=. --micro_out=.  flow_adapter/proto/adapter/adapter.proto
	protoc --proto_path=. --go_out=. --micro_out=.  telegram_bot/proto/bot_message/bot_message.proto
