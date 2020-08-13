package main

import (
	"context"
	"fmt"

	pb "github.com/matvoy/chat_server/flow_adapter/proto/adapter"
	pbtelegram "github.com/matvoy/chat_server/telegram_bot/proto/bot_message"

	"github.com/rs/zerolog"
)

type Service interface {
	SendMessageToFlow(ctx context.Context, req *pb.MessageToFlow, res *pb.Response) error
}

type flowService struct {
	log    *zerolog.Logger
	client pbtelegram.TelegramBotService
}

func NewFlowService(log *zerolog.Logger, client pbtelegram.TelegramBotService) *flowService {
	return &flowService{
		log,
		client,
	}
}

func (s *flowService) SendMessageToFlow(ctx context.Context, req *pb.MessageToFlow, res *pb.Response) error {
	s.log.Info().Msg("accept message")
	_, err := client.ProcessMessageFromFlow(ctx, &pbtelegram.MessageFromFlow{
		Text:           fmt.Sprintf("Received message: %s", req.Text),
		ExternalUserId: req.ExternalUserId,
		SessionId:      req.SessionId,
	})
	if err == nil {
		res.Success = true
	}
	return nil
}
