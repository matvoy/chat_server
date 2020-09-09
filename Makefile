.PHONY: vendor

# download vendor
vendor:
	GO111MODULE=on go mod vendor

# start all unit tests
tests:
	go test ./...

# build storage service
build-storage:
	go build -mod=vendor -o bin/webitel.chat.service.storage ./cmd/chat_storage/*.go

# build api service
build-api:
	go build -mod=vendor -o bin/webitel.chat.service.api ./chat_api/*.go

# build telegram bot service
build-telegrambot:
	go build -mod=vendor -o bin/webitel.chat.service.telegrambot ./cmd/telegram_bot/*.go

# build viber bot service
build-viberbot:
	go build -mod=vendor -o bin/webitel.chat.service.viberbot ./viber_bot/*.go

# build whatsapp bot service
build-whatsappbot:
	go build -mod=vendor -o bin/webitel.chat.service.whatsappbot ./whatsapp_bot/*.go

# build facebook bot service
build-facebookbot:
	go build -mod=vendor -o bin/webitel.chat.service.facebookbot ./facebook_bot/*.go

# build flow service
build-flowclient:
	go build -mod=vendor -o bin/webitel.chat.service.flowclient ./cmd/flow_client/*.go

# build chat server
build: build-storage build-telegrambot build-flowclient build-viberbot

# start linter
lint:
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./...

# cd & generate boiler models
generate-boiler:
	sqlboiler --wipe --no-tests -o ./models -c ./configs/sqlboiler.toml psql

proto:
	./scripts/protoc.sh

run-storage: build-storage
	./bin/webitel.chat.service.storage --registry="consul" --registry_address="consul" --store="redis" --store_table="chat:" --store_address="redis" --db_host="postgres" --db_user="postgres" --db_name="postgres" --db_password="postgres" --log_level="trace"

run-telegrambot: build-telegrambot
	./bin/webitel.chat.service.telegrambot --registry="consul" --registry_address="consul" --tg_webhook_address="example.com" --app_port=8889

run-viberbot: build-viberbot
	./bin/webitel.chat.service.viberbot --registry="consul" --registry_address="consul" --viber_webhook_address="example.com" --app_port=8889

run-whatsappbot: build-whatsappbot
	./bin/webitel.chat.service.whatsappbot --registry="consul" --registry_address="consul"

run-facebookbot: build-facebookbot
	./bin/webitel.chat.service.facebookbot --registry="consul" --registry_address="consul" --fb_webhook_address="example.com" --app_port=8889

run-flowclient: build-flowclient
	./bin/webitel.chat.service.flowclient --registry="consul" --registry_address="consul" --store="redis" --store_table="chat:" --store_address="redis" --conversation_timeout_sec=600

generate-ssl:
	openssl genrsa -out server.key 2048
	openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650