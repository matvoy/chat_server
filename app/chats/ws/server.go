package ws

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
)

type websocketChat struct {
	port int
	log  *zerolog.Logger
}

func NewWebsocketChat(port int, log *zerolog.Logger) *websocketChat {
	return &websocketChat{
		port,
		log,
	}
}

func (w websocketChat) Start() error {
	hub := newHub()
	go hub.run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	return http.ListenAndServe(fmt.Sprintf(":%v", w.port), nil)
}
