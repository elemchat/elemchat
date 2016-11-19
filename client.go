package main

import (
	"fmt"
	"os"
	"time"

	"github.com/elemchat/elemchat/codec"
	"github.com/elemchat/elemchat/conn"
	"github.com/elemchat/elemchat/msg"
	"github.com/elemchat/elemchat/wizard"
	"github.com/gorilla/websocket"
)

type ClientConfig struct {
	Addr   string
	Name   string
	Wizard *wizard.Config
}

func client(args []string) {
	cfg := &ClientConfig{
		Addr: "ws://localhost:9000/elemchat",
		Name: "mofon",
		Wizard: &wizard.Config{
			ReadTimeout:  40 * time.Second,
			WriteTimeout: 40 * time.Second,
		},
	}

	dialer := &websocket.Dialer{}

	err := match(cfg, dialer)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

func match(cfg *ClientConfig, dialer *websocket.Dialer) error {
	ws, _, err := dialer.Dial(cfg.Addr+"/match", nil)
	if err != nil {
		return err
	}

	recv := make(chan wizard.Message)
	c := conn.WebSocket(ws)

	w := wizard.NewWithConfig(cfg.Wizard, "", c, codec.JsonCodec(), recv)
	go sendPing(cfg, w)
	matchCh := make(chan *msg.Match, 1)
	go func() { // handle recv
		for message := range recv {
			switch msg.GetType(message.Msg()) {
			case msg.MATCH:
				m := message.Msg().(*msg.Match)
				select {
				case matchCh <- m:
					w.Close(false)
				default:
				}
			}
		}
	}()
	go func() { // send auth
		w.Send(&msg.Auth{UserName: cfg.Name})
	}()

	match := <-matchCh
	fmt.Println(match.FieldURI)
	fmt.Println(match.Password)
	return nil
}

func sendPing(cfg *ClientConfig, w *wizard.Wizard) {
	timeout_small := cfg.Wizard.ReadTimeout
	timeout_big := cfg.Wizard.WriteTimeout
	if timeout_big < timeout_small {
		timeout_small, timeout_big = timeout_big, timeout_small
	}
	internal := timeout_small * 4 / 3
	if internal == 0 {
		return
	}

	for {
		if !w.Closed() {
			w.Send(&msg.Ping{})
		}
		time.Sleep(internal)
	}
}
