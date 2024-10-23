package daemon

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/tera-insights/ticrypt-file-copy/copy"
)

type daemon struct {
	allowed_hosts []string
	port          string
	listener      *webSocketListener
}

func NewDaemon(port string, allowed_hosts []string) *daemon {
	daemon := &daemon{
		port:          port,
		listener:      newWebSocketListener(),
		allowed_hosts: allowed_hosts,
	}
	return daemon
}

func (d *daemon) Start() error {
	d.listener.Register("copy", func(ctx context.Context, data json.RawMessage) {
		conn := ctx.Value("conn").(*websocket.Conn)
		MsgID := ctx.Value("msg_id").(string)
		var copyMsg struct {
			SourceFilepath      string `json:"sourceFilepath"`
			DestinationFilePath string `json:"destinationFilePath"`
			ChunkSize           int    `json:"chunkSize"`
		}
		if err := json.Unmarshal(data, &copyMsg); err != nil {
			log.Printf("error unmarshalling copy message: %s\n", err.Error())
			return
		}
		progress := make(chan copy.Progress)
		go func() {
			for p := range progress {
				err := conn.WriteJSON(message{
					MsgID: MsgID,
					Event: "progress",
					Data:  json.RawMessage(fmt.Sprintf(`{"bytesWritten": %d, "totalBytes": %d}`, p.BytesWritten, p.TotalBytes)),
				})
				if err != nil {
					log.Printf("error writing progress: %s\n", err.Error())
				}
			}
		}()
		if copyMsg.SourceFilepath == "" || copyMsg.DestinationFilePath == "" {
			err := conn.WriteMessage(websocket.TextMessage, []byte("sourceFilepath and destinationFilePath are required"))
			if err != nil {
				log.Printf("error writing error message: %s\n", err.Error())
			}
			return
		}
		copier := copy.NewCopier(copyMsg.SourceFilepath, copyMsg.DestinationFilePath, copyMsg.ChunkSize, progress)
		err := copier.Copy(copy.Read, copy.Write)
		if err != nil {
			log.Printf("error copying file: %s\n", err.Error())
		}
	})

	d.listener.Register("benchmark", func(ctx context.Context, data json.RawMessage) {
		conn := ctx.Value("conn").(*websocket.Conn)
		MsgID := ctx.Value("msg_id").(string)
		var copyMsg struct {
			SourceFilepath      string `json:"sourceFilepath"`
			DestinationFilePath string `json:"destinationFilePath"`
			ChunkSize           int    `json:"chunkSize"`
		}
		if err := json.Unmarshal(data, &copyMsg); err != nil {
			log.Printf("error unmarshalling copy message: %s\n", err.Error())
			return
		}
		progress := make(chan copy.Progress)
		if copyMsg.SourceFilepath == "" {
			copyMsg.SourceFilepath = "source"
		}
		if copyMsg.DestinationFilePath == "" {
			copyMsg.DestinationFilePath = "destination"
		}
		copier := copy.NewCopier(copyMsg.SourceFilepath, copyMsg.DestinationFilePath, copyMsg.ChunkSize, progress)

		if err := conn.WriteMessage((websocket.TextMessage), []byte("Starting Benchmark")); err != nil {
			log.Printf("error writing start benchmark message: %s\n", err.Error())
		}

		go func() {
			for stat := range progress {
				if err := conn.WriteJSON(message{
					MsgID: MsgID,
					Event: "progress",
					Data:  json.RawMessage(fmt.Sprintf(`{"bytesWritten": %d, "totalBytes": %d}`, stat.BytesWritten, stat.TotalBytes)),
				}); err != nil {
					log.Printf("error writing progress: %s\n", err.Error())
				}
			}
		}()

		if err := copier.Benchmark(copy.Read, copy.Write); err != nil {
			log.Printf("error benchmarking file: %s\n", err.Error())
			err := conn.WriteMessage((websocket.TextMessage), []byte(fmt.Sprintf("error benchmarking file: %s\n", err.Error())))
			if err != nil {
				log.Printf("error writing error message: %s\n", err.Error())
			}
			return
		}

		if err := conn.WriteMessage((websocket.TextMessage), []byte("Benchmark Complete")); err != nil {
			log.Printf("error writing benchmark complete message: %s\n", err.Error())
		}
	})

	d.listener.Register("ping", func(ctx context.Context, data json.RawMessage) {
		conn := ctx.Value("conn").(*websocket.Conn)
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Printf("error writing pong: %s\n", err.Error())
		}
	})

	d.listener.Register("stop", func(ctx context.Context, data json.RawMessage) {
		conn := ctx.Value("conn").(*websocket.Conn)
		d.Close()
		err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Printf("error writing close message: %s\n", err.Error())
		}
	})

	http.HandleFunc("/ws", d.Serve)
	addr := flag.String("addr", "localhost:4242", "http service address")
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fmt.Println("Shutting Down Deamon")
	return nil
}

func (d *daemon) Close() {
	d.listener.Close()
}

var upgrader = websocket.Upgrader{}

func (d *daemon) Serve(w http.ResponseWriter, r *http.Request) {

	if !isHostAllowed(r.RemoteAddr, d.allowed_hosts) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		fmt.Printf("Forbidden request from %s, only %s allowed \n", r.Host, d.allowed_hosts)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer ws.Close()

	fmt.Printf("Listening for messages on connection to: %s \n", ws.RemoteAddr())
	err = d.listener.Listen(ws)
	if err != nil {
		log.Println("listen:", err)
		return
	}
	fmt.Printf("Connection to %s closed \n", ws.RemoteAddr())
}

func isHostAllowed(host string, allowed_hosts []string) bool {
	for _, allowed_host := range allowed_hosts {
		if strings.Split(host, ":")[0] == allowed_host {
			return true
		}
	}
	return false
}
