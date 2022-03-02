package websocket

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type SocketPayload struct {
	Client  *Client     `json:"-"`
	Type    string      `json:"type,omitempty"`
	Message interface{} `json:"message,omitempty"`
}

type ConnectionDetail struct {
	UserID          string
	IPAddress       string
	IsSubscribing   bool
	AssetSubscribed string
}

type Client struct {
	Conn *websocket.Conn
	Pool *Pool
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		payload := SocketPayload{}
		err := c.Conn.ReadJSON(&payload)
		if err != nil {
			log.Println(err)
			return
		}
		payload.Client = c

		c.Pool.Message <- payload
		fmt.Printf("Message Received: %+v\n", payload)
	}
}
