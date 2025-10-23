package websocket

import (
	"log"
	"net/http"
	"slices"

	"github.com/MilindGour/jellyfin-media-renamer/filesystem"
	ws "github.com/gorilla/websocket"
)

type JMRWebSocket struct {
	connections []*ws.Conn
	upgrader    *ws.Upgrader
}

func NewJMRWebSocket() JMRWebSocket {
	return JMRWebSocket{
		connections: []*ws.Conn{},
		upgrader: &ws.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (j *JMRWebSocket) UpgradeConnectionAndAddClient(w http.ResponseWriter, r *http.Request) error {
	conn, err := j.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	j.AddConnection(conn)

	return nil
}

func (j *JMRWebSocket) AddConnection(conn *ws.Conn) {
	log.Println("=== A new client is being connected via WebSocket ===")
	j.connections = append(j.connections, conn)
}

func (j *JMRWebSocket) RemoveConnection(conn *ws.Conn) {
	index := slices.Index(j.connections, conn)
	conn.Close()

	if index > -1 {
		j.connections = append(j.connections[:index], j.connections[index+1:]...)
	}
}

func (j *JMRWebSocket) SendMessage(message any) {
	if len(j.connections) > 0 {
		for _, conn := range j.connections {
			conn.WriteJSON(message)
		}
	}
}

func (j *JMRWebSocket) SendWSEvent(eventData EventData) {
	j.SendMessage(EventMessage{
		Message: Message{
			Type: "event",
		},
		Data: eventData,
	})
}

func (j *JMRWebSocket) SendProgressMessage(progress []filesystem.FileTransferProgress) {
	j.SendMessage(ProgressMessage{
		Message: Message{
			Type: "progress",
		},
		Data: progress,
	})
}

type Message struct {
	Type string `json:"type"`
}

type ProgressMessage struct {
	Message
	Data []filesystem.FileTransferProgress `json:"data"`
}

type EventData struct {
	Name   string `json:"name"`
	Detail any    `json:"detail"`
}

type EventMessage struct {
	Message
	Data EventData `json:"data"`
}
