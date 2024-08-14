package daemon

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	ticrypt "github.com/tera-insights/ticrypt-go"
)

type webSocketListener struct {
	hostID   string
	tcClient ticrypt.Client

	conn     *websocket.Conn
	stop     chan any
	handlers map[string]func(json.RawMessage)
}

func newWebSocketListener(hostID string, tcClient ticrypt.Client) *webSocketListener {
	listener := &webSocketListener{
		hostID:   hostID,
		tcClient: tcClient,
		stop:     make(chan any),
		handlers: make(map[string]func(json.RawMessage)),
	}
	return listener
}

func (l *webSocketListener) Register(updateType string, f func(json.RawMessage)) {
	l.handlers[updateType] = f
}

func (l *webSocketListener) Unregister(updateType string) {
	delete(l.handlers, updateType)
}

func (l *webSocketListener) Close() {
	close(l.stop)
}

func (l *webSocketListener) run() {
	for {
		select {
		case <-l.stop:
			if l.conn != nil {
				l.conn.Close()
			}
			return
		default:
			conn, err := l.GetConnection()
			if err != nil {
				log.Printf("error getting connection: %s\n", err.Error())
				log.Printf("retrying in 5 seconds\n")
				time.Sleep(5 * time.Second)
				continue
			}

			msgType, data, err := conn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
					log.Printf("error reading message: %s\n", err.Error())
				}

				conn.Close()
				l.conn = nil
				continue
			}

			if msgType != websocket.TextMessage {
				log.Printf("unexpected message of type %d: %+v\n", msgType, data)
				continue
			}

			var msg struct {
				Type  string          `json:"type"`
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

			go handler(msg.Data)
		}
	}
}

func (l *webSocketListener) GetConnection() (*websocket.Conn, error) {
	if l.conn == nil {
		conn, err := l.tcClient.NewConnToHostUpdates(l.hostID)
		if err != nil {
			return nil, err
		}

		log.Printf("connected to host updates\n")
		l.conn = conn
	}

	return l.conn, nil
}
