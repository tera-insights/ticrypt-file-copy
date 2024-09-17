package daemon

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type webSocketListener struct {
	stop     chan any
	handlers map[string]func(*websocket.Conn, json.RawMessage)
}

func newWebSocketListener() *webSocketListener {
	listener := &webSocketListener{
		stop:     make(chan any),
		handlers: make(map[string]func(*websocket.Conn, json.RawMessage)),
	}
	return listener
}

func (l *webSocketListener) Register(updateType string, f func(*websocket.Conn, json.RawMessage)) {
	l.handlers[updateType] = f
}

func (l *webSocketListener) Unregister(updateType string) {
	delete(l.handlers, updateType)
}

func (l *webSocketListener) Close() {
	close(l.stop)
}

func (l *webSocketListener) Listen(conn *websocket.Conn) error {
	defer conn.Close()
	for {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Printf("error reading message: %s\n", err.Error())
				return err
			}
			return nil
		}

		if msgType != websocket.TextMessage {
			log.Printf("unexpected message of type %d: %+v\n", msgType, data)
			continue
		}

		var msg struct {
			MsgID string          `json:"msg_id"`
			Event string          `json:"event"`
			Data  json.RawMessage `json:"data"`
		}
		err = json.Unmarshal(data, &msg)
		if err != nil {
			log.Printf("error unmarshalling message: %s\n", err.Error())
			continue
		}

		handler, ok := l.handlers[msg.Event]
		if !ok {
			log.Printf("no handler for event: %s\n", msg.Event)
			continue
		}

		go handler(conn, msg.Data)
	}
}
