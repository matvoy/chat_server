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
build-flowadapter:
	go build -mod=vendor -o bin/webitel.chat.service.flowadapter ./flow_adapter/*.go

# build chat server
build: build-storage build-telegrambot build-flowadapter

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

run-storage: build-storage
	./bin/webitel.chat.service.storage --db_host="postgres" --db_user="postgres" --db_name="postgres" --db_password="postgres" --log_level="trace"

run-telegrambot: build-telegrambot
	./bin/webitel.chat.service.telegrambot --store="redis" --store_table="chat:" --store_address="redis" --telegram_bot_token="token" --profile_id=1 --conversation_timeout=300

run-flowadapter: build-flowadapter
	./bin/webitel.chat.service.flowadapter

run-api:
	./bin/webitel.chat.service.api start  --api_http_port=55020 --db_host="postgres" --db_user="postgres" --db_name="postgres" --db_password="postgres"

run: run-storage run-telegrambot run-flowadapter