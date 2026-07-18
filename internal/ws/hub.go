package ws

import (
	"net/http"
	"realtime-log-aggregation/internal/models"

	"github.com/gorilla/websocket"
)

type Hub struct{
	Clients	    map[*websocket.Conn]bool
	Broadcast   chan models.LogEvent
	Register    chan *websocket.Conn
	Unregister  chan *websocket.Conn
}

var Upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {return true},
}

func NewHub (broadcast chan models.LogEvent) *Hub {
	var newHub = Hub{
		Clients: make(map[*websocket.Conn]bool),
		Broadcast: broadcast,
		Register: make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
	}
	return &newHub
}

func (h *Hub) Run(){
	// for {
	// 	select {
	// 	case Register: 
	// 		h.Clients.append(h.Register)
	// 	case Unregister:
	// 		if h.Clients {
	// 			delete(h.Clients[h.Unregister], true)
	// 			h.Clients[h.Unregister].close()
	// 		}
	// 	case Broadcast:
	// 		for client := range h.Clients{

	// 		} 
	// 	}
	// }
}

func (h *Hub) HandleWS (w *http.ResponseWriter, r *http.Request) error {
	// upgrade http server connection to websocket protocol
	// websocketConn, connErr := Upgrader.Upgrade(w, r, nil)
	// if connErr != nil {
	// 	return fmt.Errorf("Http connection upgrade failed, %v", connErr)
	// }
	// h.Register <- websocketConn
	// go func ()  {
		
	// }()
	return nil
}

