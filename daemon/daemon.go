package daemon

import (
	"encoding/json"
	"log"

	"github.com/tera-insights/ticrypt-file-copy/copy"
	ticrypt "github.com/tera-insights/ticrypt-go"
)

type daemon struct {
	tcClient *ticrypt.Client
	listener *webSocketListener
	hostID   string
}

func NewDaemon(hostID string, tcClient *ticrypt.Client) *daemon {
	daemon := &daemon{
		tcClient: tcClient,
		listener: newWebSocketListener(hostID, *tcClient),
		hostID:   hostID,
	}
	return daemon
}

func (d *daemon) Start() {
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
		copier := copy.NewCopier(copyMsg.SourceFilepath, copyMsg.DestinationFilePath, copyMsg.ChunkSize)
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
			copier := copy.NewCopier(copyMsg.SourceFilepath, copyMsg.DestinationFilePath, copyMsg.ChunkSize)
			copier.Benchmark(copy.Read, copy.Write)
			return
		}
		return
	})
	d.listener.Register("stop", func(data json.RawMessage) {
		d.Close()
		return
	})

	go d.listener.run()
}

func (d *daemon) Close() {
	d.listener.Close()
}
