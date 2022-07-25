package ws

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type WSClient struct {
	hub *WSHub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	wg sync.WaitGroup
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (c *WSClient) readPump() {
	log.Println("WSClient: reading")
	defer func() {
		log.Println("WSClient: readPump finish")
		c.wg.Wait()
		log.Println("WSClient: unregister")
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		log.Println("WSClient RUN")
		if c.hub.closed {
			log.Println("WSClient hub closed")
			break
		}
		var e WSEvent
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		err = json.Unmarshal(message, &e)
		if err == nil {
			log.Println(e)
			if e.Event == "init" {
				var init WSEventInit
				json.Unmarshal(message, &init)
				if err == nil {
					go func() {
						defer func() {
							c.wg.Done()
							log.Println("WSClient WG DONE")
						}()
						log.Println("Handle", init)
						time.Sleep(time.Second * 10)
						log.Println("Handle finis")
					}()
					c.wg.Add(1)
					log.Println("WSClient WG ADD")
				}
			}
		}

		//c.hub.broadcast <- "hello"
	}
}

type WSEvent struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type WSEventInit struct {
	Event string          `json:"event"`
	Data  WSEventInitData `json:"data"`
}

type WSEventInitData struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
}

func (c *WSClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
