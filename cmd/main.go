package main

import (
	"log"
	"multiclient-websocket/internal/boot"
)

func main() {
	if err := boot.Server(); err != nil {
		log.Println("[Server] failed to boot websocket server due to " + err.Error())
	}
}
