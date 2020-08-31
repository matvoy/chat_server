.PHONY: vendor

# download vendor
vendor:
	GO111MODULE=on go mod vendor

# start all unit tests
tests:
	go test ./...

# build storage service
build-storage:
	go build -mod=vendor -o bin/webitel.chat.service.storage ./chat_storage/*.go

# build api service
build-api:
	go build -mod=vendor -o bin/webitel.chat.service.api ./chat_api/*.go

# build telegram bot service
build-telegrambot:
	go build -mod=vendor -o bin/webitel.chat.service.telegrambot ./telegram_bot/*.go

# build telegram bot service
build-flowclient:
	go build -mod=vendor -o bin/webitel.chat.service.flowclient ./flow_client/*.go

# build chat server
build: build-storage build-telegrambot build-flowclient

# start linter
lint:
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./...

# cd & generate boiler models
generate-boiler:
	sqlboiler --wipe --no-tests -o ./models -c ./configs/sqlboiler.toml psql

proto: 
	protoc --proto_path=. --go_out=. --micro_out=.  chat_storage/proto/storage/storage.proto
	protoc --proto_path=. --go_out=. --micro_out=.  flow_client/proto/flow_client/flow_client.proto
	protoc --proto_path=. --go_out=. --micro_out=.  flow_client/proto/flow_manager/flow_manager.proto
	protoc --proto_path=. --go_out=. --micro_out=.  telegram_bot/proto/bot_message/bot_message.proto

run-storage: build-storage
	./bin/webitel.chat.service.storage --registry="consul" --registry_address="consul" --store="redis" --store_table="chat:" --store_address="redis" --db_host="postgres" --db_user="postgres" --db_name="postgres" --db_password="postgres" --log_level="trace"

run-telegrambot: build-telegrambot
	./bin/webitel.chat.service.telegrambot --registry="consul" --registry_address="consul" --tg_webhook_address="example.com"

run-flowclient: build-flowclient
	./bin/webitel.chat.service.flowclient --registry="consul" --registry_address="consul" --store="redis" --store_table="chat:" --store_address="redis" --conversation_timeout_sec=600

generate-ssl:
	openssl genrsa -out server.key 2048
	openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650