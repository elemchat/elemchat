package main

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/elemchat/elemchat/codec"
	"github.com/elemchat/elemchat/conn"
	"github.com/elemchat/elemchat/field"
	"github.com/elemchat/elemchat/msg"
	"github.com/elemchat/elemchat/wizard"
	"github.com/elemchat/elemchat/wizard/attr"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type ServerConfig struct {
	WebSocketAddr string
	SSL           bool
	Wizard        wizard.Config
}

var MatchQueue chan *wizard.Wizard

func server(args []string) {
	cfg := &ServerConfig{
		WebSocketAddr: "localhost:9000",
		SSL:           false,
		Wizard: wizard.Config{
			ReadTimeout:  time.Second * 40,
			WriteTimeout: time.Second * 40,
			DefaultAttr: attr.Attr{
				Blood: 20,
			},
		},
	}

	log := logrus.StandardLogger()
	log.Level = logrus.DebugLevel

	router := mux.NewRouter().PathPrefix("/elemchat/").
		Subrouter().StrictSlash(true)
	registerMatch(log, router, cfg)

	server := &http.Server{
		Addr:         cfg.WebSocketAddr,
		Handler:      router,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}

	log.WithField("addr", server.Addr).Info("server start")
	server.ListenAndServe()
}

func registerMatch(log *logrus.Logger, router *mux.Router, cfg *ServerConfig) {
	var match *field.Field
	match = field.New(func(message wizard.Message) {

		log.WithField("type", msg.GetType(message.Msg())).
			Debug("handle message")

		switch msg.GetType(message.Msg()) {
		case msg.AUTH:
			m := message.Msg().(*msg.Auth)
			if m.UserName == "" {
				return
			}

			match.WithLock(func(match *field.Field) {
				if message.Wizard().Name != "" {
					return
				}

				message.Wizard().Name = m.UserName
				log.WithField("username", message.Wizard().Name).
					Info("match auth")
			})
		case msg.PING:
			match.WithLock(func(match *field.Field) {
				go message.Wizard().Send(&msg.Pong{})
			})
		}
	})
	router.Path("/match").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ws, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				return
			}

			c := conn.WebSocket(ws)
			match.Enter(func(recv chan<- wizard.Message) *wizard.Wizard {
				return wizard.NewWithConfig(&cfg.Wizard,
					"", c, codec.JsonCodec(), recv)
			})
		})

}
