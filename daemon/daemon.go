package daemon

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tera-insights/ticrypt-file-copy/copy"
)

type daemon struct {
	allowed_hosts []string
	port          string
	listener      *webSocketListener
}

type ticpConfig struct {
	Host string `toml:"host"`
}

func NewDaemon(port string, allowed_hosts []string) *daemon {
	daemon := &daemon{
		port:     port,
		listener: newWebSocketListener(),
	}
	return daemon
}

func (d *daemon) Start() error {
	d.listener.Register("copy", func(data json.RawMessage) {
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
		copier := copy.NewCopier(copyMsg.SourceFilepath, copyMsg.DestinationFilePath, copyMsg.ChunkSize, progress)
		copier.Copy(copy.Read, copy.Write)
		return
	})
	d.listener.Register("benchmark", func(data json.RawMessage) {
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
		copier := copy.NewCopier(copyMsg.SourceFilepath, copyMsg.DestinationFilePath, copyMsg.ChunkSize, progress)
		copier.Benchmark(copy.Read, copy.Write)
		return
	})
	d.listener.Register("stop", func(data json.RawMessage) {
		d.Close()
		return
	})

	http.HandleFunc("/ws", d.Serve)
	http.ListenAndServe(fmt.Sprintf(":%s", d.port), nil)
	return nil
}

func (d *daemon) Close() {
	d.listener.Close()
}

var upgrader = websocket.Upgrader{}

func internalError(ws *websocket.Conn, msg string, err error) {
	log.Println(msg, err)
	ws.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
}

func (d *daemon) Serve(w http.ResponseWriter, r *http.Request) {

	if !isHostAllowed(r.Host, d.allowed_hosts) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	d.listener.Listen(ws)

	defer ws.Close()
}

func isHostAllowed(host string, allowed_hosts []string) bool {
	for _, allowed_host := range allowed_hosts {
		if host == allowed_host {
			return true
		}
	}
	return false
}
