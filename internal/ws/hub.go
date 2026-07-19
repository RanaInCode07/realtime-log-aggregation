package ws

import (
	"log"
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

func (h *Hub) Run() {
	for {
		select {
		case client:= <-h.Register: 
			h.Clients[client] = true
		case client:= <-h.Unregister:
			if h.Clients[client] {
				delete(h.Clients, client)
				client.Close()
			}
		case logItem:= <-h.Broadcast:
			for client := range h.Clients{
				if writeJsonErr := client.WriteJSON(logItem); writeJsonErr != nil{
					// go routine to avoid deadlock as channel are unbuffered so it will stuck forever until someone read it
					// since the hub itself was the only one that could read it and it is current stuck standing in line to send it
					go func(c *websocket.Conn){
						h.Unregister <- c
					}(client)
				}
			} 
		}
	}
}

func (h *Hub) HandleWS (w http.ResponseWriter, r *http.Request) {
	//upgrade http server connection to websocket protocol
	websocketConn, connErr := Upgrader.Upgrade(w, r, nil)
	if connErr != nil {
		log.Printf("Error during upgrading http connection to websocket protocol, %v \n", connErr)
		return
	}
	h.Register <- websocketConn
	go func ()  {
		for {
			_, _, err := websocketConn.ReadMessage()
		    if err != nil {
			h.Unregister <- websocketConn
			break
		}
		}
	}()
}

