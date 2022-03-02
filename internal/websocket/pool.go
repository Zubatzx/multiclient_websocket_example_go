package websocket

import (
	"fmt"
	"strings"
	"time"
)

type Pool struct {
	// Client Related
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]ConnectionDetail
	// Etc
	Message chan SocketPayload
}

func NewPool() *Pool {
	return &Pool{
		// Client Related
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]ConnectionDetail),
		// Etc
		Message: make(chan SocketPayload),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = ConnectionDetail{}
			fmt.Println("New Connection Registered")
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Connection Unregistered")
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
		case message := <-pool.Message:
			fmt.Println("New Message Received")
			pool.handleIO(message)
		}
	}
}

func (pool *Pool) Publish() {
	for {
		time.Sleep(time.Second * 10)
		fmt.Println("Publish message to all clients in Pool", len(pool.Clients))
		for client, detail := range pool.Clients {
			if detail.IsSubscribing {
				var message string
				switch strings.ToLower(detail.AssetSubscribed) {
				case "bbca":
					message = "Normal"
				case "buka":
					message = "To The Moon, BUT REVERSED!"
				case "tkpd":
					message = "To The Moon!"
				default:
					continue
				}
				fmt.Println("Publish asset detail to", &client)
				if err := client.Conn.WriteJSON(SocketPayload{Message: message}); err != nil {
					fmt.Println(err)
				}

			}
		}
	}
}

func (pool *Pool) handleIO(message SocketPayload) {
	switch message.Type {
	case Subscribe:
		fmt.Println("Client Subscribe")
		pool.Clients[message.Client] = ConnectionDetail{
			IsSubscribing:   true,
			AssetSubscribed: message.Message.(string),
		}
	case Unsubscribe:
		fmt.Println("Client Unsubscribe")
		pool.Clients[message.Client] = ConnectionDetail{
			IsSubscribing:   false,
			AssetSubscribed: "",
		}
	default:
		fmt.Println("Unavailable Type")
	}
}
