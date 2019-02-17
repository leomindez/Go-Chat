package main

import (
	"net/http"

	"github.com/leomindez/ChatGo/pkg/chatwebsocket"
)

const port = ":8080"

func routes() {
	chatwebsocket.HandleWebSocketConnection()
}

func startServer(port string, handler http.Handler) {
	http.ListenAndServe(port, handler)
}

func main() {
	routes()
	startServer(port, nil)
}
