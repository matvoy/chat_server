#!/bin/sh
set -x

protoc -I api/proto/entity --go_out=api/proto/entity --micro_out=api/proto/entity api/proto/entity/entity.proto
mv ./api/proto/entity/github.com/matvoy/chat_server/api/proto/entity/* ./api/proto/entity/
rm -rf ./api/proto/entity/github.com

protoc -I api/proto/entity -I api/proto/chat_storage --go_out=api/proto/chat_storage --micro_out=api/proto/chat_storage api/proto/chat_storage/chat_storage.proto
mv ./api/proto/chat_storage/github.com/matvoy/chat_server/api/proto/chat_storage/* ./api/proto/chat_storage/
rm -rf ./api/proto/chat_storage/github.com

protoc -I api/proto/entity -I api/proto/flow_client --go_out=api/proto/flow_client --micro_out=api/proto/flow_client api/proto/flow_client/flow_client.proto
mv ./api/proto/flow_client/github.com/matvoy/chat_server/api/proto/flow_client/* ./api/proto/flow_client/
rm -rf ./api/proto/flow_client/github.com

protoc -I api/proto/entity -I api/proto/flow_manager --go_out=api/proto/flow_manager --micro_out=api/proto/flow_manager api/proto/flow_manager/flow_manager.proto
mv ./api/proto/flow_manager/github.com/matvoy/chat_server/api/proto/flow_manager/* ./api/proto/flow_manager/
rm -rf ./api/proto/flow_manager/github.com

protoc -I api/proto/entity -I api/proto/telegram_bot --go_out=api/proto/telegram_bot --micro_out=api/proto/telegram_bot api/proto/telegram_bot/telegram_bot.proto
mv ./api/proto/telegram_bot/github.com/matvoy/chat_server/api/proto/telegram_bot/* ./api/proto/telegram_bot/
rm -rf ./api/proto/telegram_bot/github.com
