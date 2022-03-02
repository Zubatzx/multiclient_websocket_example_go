package boot

import (
	"fmt"
	"log"
	"net/http"

	"multiclient-websocket/internal/config"

	websocketS "multiclient-websocket/internal/websocket"
	"multiclient-websocket/pkg/websocket"
)

func Server() error {
	// Initialize context
	// ctx := context.Background()

	// Initialize config from file
	err := config.Init()
	if err != nil {
		log.Fatalf("[CONFIG] Failed to initialize config: %v", err)
	}
	cfg := config.Get()

	// Initialized pool of channels
	pool := websocketS.NewPool()
	go pool.Start()
	go pool.Publish()

	// Initialized Server and Client Websocket Connection
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("WebSocket Endpoint Hit")
		conn, err := websocket.Upgrade(w, r)
		if err != nil {
			fmt.Fprintf(w, "%+v\n", err)
		}

		client := &websocketS.Client{
			Conn: conn,
			Pool: pool,
		}

		pool.Register <- client
		client.Read()
	})

	fmt.Println("Server starting at", cfg.Server.Port)
	http.ListenAndServe(cfg.Server.Port, nil)

	return nil
}
