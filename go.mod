module github.com/matvoy/chat_server

go 1.14

replace github.com/nats-io/nats.go v1.8.2-0.20190607221125-9f4d16fe7c2d => github.com/nats-io/nats.go v1.8.1

require (
	github.com/friendsofgo/errors v0.9.2
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d // indirect
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.2
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/lib/pq v1.8.0
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/broker/rabbitmq/v2 v2.9.1
	github.com/micro/go-plugins/registry/consul/v2 v2.9.1
	github.com/micro/go-plugins/store/redis/v2 v2.9.1
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/rs/zerolog v1.19.0
	github.com/volatiletech/null/v8 v8.1.0
	github.com/volatiletech/sqlboiler/v4 v4.2.0
	github.com/volatiletech/strmangle v0.0.1
	github.com/webitel/protos v0.0.0-20201009130933-65f553ce8ca3
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899 // indirect
	google.golang.org/grpc v1.33.0 // indirect
	google.golang.org/grpc/examples v0.0.0-20201008225051-06c094c3ab22 // indirect
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace github.com/webitel/protos => ../protos
