package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

const HOSTNAME = "localhost"
const PORT = "4242"

func main() {
	addr := flag.String("addr", fmt.Sprintf("%s:%s", HOSTNAME, PORT), "http service address")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	data, err := json.Marshal(
		struct {
			MsgID string          `json:"msg_id"`
			Event string          `json:"event"`
			Data  json.RawMessage `json:"data"`
		}{
			MsgID: "1",
			Event: "benchmark",
			Data:  json.RawMessage(`{"benchmark": "benchmark"}`),
		},
	)
	if err != nil {
		log.Println("json marshal:", err)
		return
	}

	err = c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Println("write:", err)
		return
	}

	done := make(chan struct{})

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				done <- struct{}{}
				return
			}
			if string(message) == "Benchmark Complete" {
				log.Println("Benchmark Complete")
				done <- struct{}{}
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	<-done
	err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		return
	}
}
