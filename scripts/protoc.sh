#!/bin/sh
set -x

protoc -I api/proto/chat --go_out=api/proto/chat --micro_out=api/proto/chat api/proto/chat/chat.proto
mv ./api/proto/chat/github.com/matvoy/chat_server/api/proto/chat/* ./api/proto/chat/
rm -rf ./api/proto/chat/github.com

protoc -I api/proto/flow_manager --go_out=api/proto/flow_manager --micro_out=api/proto/flow_manager api/proto/flow_manager/flow_manager.proto
# mv ./api/proto/flow_manager/github.com/matvoy/chat_server/api/proto/flow_manager/* ./api/proto/flow_manager/
# rm -rf ./api/proto/flow_manager/github.com

protoc -I api/proto/chat -I api/proto/bot --go_out=api/proto/bot --micro_out=api/proto/bot api/proto/bot/bot.proto
mv ./api/proto/bot/github.com/matvoy/chat_server/api/proto/bot/* ./api/proto/bot/
rm -rf ./api/proto/bot/github.com

protoc -I api/proto/auth --go_out=api/proto/auth --micro_out=api/proto/auth api/proto/auth/authN.proto