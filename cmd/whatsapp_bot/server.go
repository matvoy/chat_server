package main

import (
	"context"
	"time"

	pbstorage "github.com/matvoy/chat_server/cmd/chat_storage/proto/storage"
	pb "github.com/matvoy/chat_server/cmd/whatsapp_bot/proto/bot_message"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	whatsapp "github.com/Rhymen/go-whatsapp"
	"github.com/rs/zerolog"
)

type ChatServer interface {
	MessageFromFlow(ctx context.Context, req *pb.MessageFromFlowRequest, res *pb.MessageFromFlowResponse) error
	StartServer() error
}

type whatsappBot struct {
	log       *zerolog.Logger
	client    pbstorage.StorageService
	bot       *whatsapp.Conn
	startTime int64
}

func NewWhatsappBotServer(
	log *zerolog.Logger,
	client pbstorage.StorageService,
) *whatsappBot {
	b := &whatsappBot{
		log:       log,
		client:    client,
		startTime: time.Now().Unix(),
	}
	bot, err := whatsapp.NewConn(120 * time.Second)
	if err != nil {
		b.log.Error().Msg(err.Error())
	}
	bot.SetClientVersion(2, 2035, 14)
	bot.AddHandler(b)
	b.bot = bot
	return b
}

func (b *whatsappBot) StartServer() error {
	qrChan := make(chan string)
	go func() {
		terminal := qrcodeTerminal.New()
		terminal.Get(<-qrChan).Print()
	}()
	_, err := b.bot.Login(qrChan)
	return err
}

func (b *whatsappBot) HandleError(err error) {
	b.log.Error().Msg(err.Error())
}

func (b *whatsappBot) HandleTextMessage(message whatsapp.TextMessage) {

	if message.Info.FromMe || int64(message.Info.Timestamp) < b.startTime {
		return
	}
	b.log.Debug().
		Str("id", message.Info.RemoteJid).
		Str("username", message.Info.RemoteJid).
		Str("text", message.Text).
		Msg("receive message")

	msg := &pbstorage.ProcessMessageRequest{
		SessionId:      message.Info.RemoteJid,
		ExternalUserId: message.Info.RemoteJid,
		Username:       message.Info.RemoteJid,
		Text:           message.Text,
		ProfileId:      4,
	}

	res, err := b.client.ProcessMessage(context.Background(), msg)
	if err != nil || res == nil {
		b.log.Error().Msg(err.Error())
	}

}

func (b *whatsappBot) MessageFromFlow(ctx context.Context, req *pb.MessageFromFlowRequest, res *pb.MessageFromFlowResponse) error {
	text := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: req.SessionId,
		},
		Text: req.GetMessage().GetTextMessage().GetText(),
	}
	if _, err := b.bot.Send(text); err != nil {
		b.log.Error().Msg(err.Error())
	}
	return nil
}
