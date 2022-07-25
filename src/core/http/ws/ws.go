package ws

import (
	"context"

	"log"
	"net/http"
	"sync"
)

type WSHub struct {
	clients map[*WSClient]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *WSClient

	// Unregister requests from clients.
	unregister chan *WSClient

	wg sync.WaitGroup

	closed bool
}

type WSRoom struct {
	data    chan []byte
	members map[*WSClient]bool
}

func (h *WSHub) Run(ctx context.Context, wg *sync.WaitGroup) error {
	wsShutdownCh := make(chan struct{})

	log.Println("WSHub run")
	wg.Add(1)
	h.closed = false

	go func() {
		<-wsShutdownCh
		for client := range h.clients {
			client.conn.Close()
		}
		h.closed = true
		log.Println("WSHub wait clients")
		h.wg.Wait()
		log.Println("WSHub shutdown")
		wg.Done()
	}()

	go func() {
		<-ctx.Done()
		wsShutdownCh <- struct{}{}
	}()

	go func() {
		for {
			select {
			case client := <-h.register:
				log.Println("WSHub add client")
				h.wg.Add(1)
				h.clients[client] = true
			case client := <-h.unregister:
				log.Println("WSHub try client")
				if _, ok := h.clients[client]; ok {
					h.wg.Done()
					delete(h.clients, client)
					close(client.send)
					log.Println("WSHub removed client")
				}
			case message := <-h.broadcast:
				for client := range h.clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}()

	return nil
}

func (h *WSHub) Shutdown() error {
	return nil
}

func NewWSHub() *WSHub {
	return &WSHub{
		broadcast:  make(chan []byte),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
		clients:    make(map[*WSClient]bool),
	}
}

func (hub *WSHub) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &WSClient{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
