.PHONY: vendor

# download vendor
vendor:
	GO111MODULE=on go mod vendor

# start all unit tests
tests:
	go test ./...

# build chat service
build-chat:
	go build -mod=vendor -o bin/webitel.chat.server ./cmd/chat/*.go

# build bot service
build-bot:
	go build -mod=vendor -o bin/webitel.chat.bot ./cmd/bot/*.go

# build all servises
build: build-chat build-bot

# generate boiler models
generate-boiler:
	sqlboiler --wipe --no-tests -o ./models -c ./configs/sqlboiler.toml psql

proto:
	./scripts/protoc.sh

run-chat: build-chat
	./bin/webitel.chat.server --client_retries=0 \
	--registry="consul" \
	--registry_address="consul" \
	--store="redis" \
	--store_table="chat:" \
	--store_address="redis:6379" \
	--broker="rabbitmq" \
	--broker_address="amqp://user:password@rabbitmq:5672/" \
	--db_host="postgres" \
	--db_user="postgres" \
	--db_name="postgres" \
	--db_password="postgres" \
	--log_level="trace" \
	--conversation_timeout_sec=600

run-bot: build-bot
	./bin/webitel.chat.bot --client_retries=0 \
	--registry="consul" \
	--registry_address="consul" \
	--webhook_address="https://example.com/" \
	--app_port=8889

generate-ssl:
	openssl genrsa -out server.key 2048
	openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650