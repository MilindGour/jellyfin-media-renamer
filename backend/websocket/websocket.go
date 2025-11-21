package websocket

import (
	"log"
	"net/http"

	"github.com/MilindGour/jellyfin-media-renamer/filesystem"
	ws "github.com/gorilla/websocket"
)

type JMRWebSocket struct {
	connections map[string]*ws.Conn
	upgrader    *ws.Upgrader
}

func NewJMRWebSocket() JMRWebSocket {
	return JMRWebSocket{
		connections: map[string]*ws.Conn{},
		upgrader: &ws.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (j *JMRWebSocket) UpgradeConnectionAndAddClient(w http.ResponseWriter, r *http.Request, id string) error {
	conn, err := j.upgrader.Upgrade(w, r, nil)
	conn.EnableWriteCompression(true)
	if err != nil {
		return err
	}
	j.AddConnection(conn, id)

	return nil
}

func (j *JMRWebSocket) AddConnection(conn *ws.Conn, id string) {
	log.Printf("[WebSocket] add connection: %s", id)
	j.RemoveAllConnections()
	j.connections[id] = conn
	log.Printf("[WebSocket] total connections now: %d", len(j.connections))
}

func (j *JMRWebSocket) RemoveAllConnections() {
	if len(j.connections) > 0 {
		for id, c := range j.connections {
			c.Close()
			delete(j.connections, id)
		}
	}
}

func (j *JMRWebSocket) RemoveConnection(id string) {
	log.Printf("[WebSocket] remove connection: %s", id)
	theConn, isthere := j.connections[id]

	if isthere {
		theConn.Close()
		delete(j.connections, id)
	}
	log.Printf("[WebSocket] total connections now: %d", len(j.connections))
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
