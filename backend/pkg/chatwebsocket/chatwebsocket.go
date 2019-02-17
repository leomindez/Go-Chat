package chatwebsocket

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
)

const prefixWs = "/ws"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleWebSocketConnection() {
	pool := NewPool()
	go pool.Start()
	http.HandleFunc(prefixWs, func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = connectionCheckOrigin
		logMessage(r.Host)
		conn := upgradeConnection(w, r)
		listenMultipleConnections(conn, pool)
	})
}

func listenMultipleConnections(conn *websocket.Conn, pool *Pool) {

	client := &Client{
		ID:   string(rand.Intn(100)),
		Conn: conn,
		Pool: pool,
	}
	pool.Register <- client
	client.Read()
}

func readConnection(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		printError(err)
		logMessage(string(p))
		writeMessage(conn, messageType, p)
	}
}

func connectionCheckOrigin(r *http.Request) bool {
	return true
}

func writeMessage(conn *websocket.Conn, messageType int, p []byte) {
	err := conn.WriteMessage(messageType, p)
	printError(err)
}

func nextWriteMessage(conn *websocket.Conn) {
	for {
		logMessage("Sending")
		messageType, r, err := conn.NextReader()
		printError(err)
		w, err := conn.NextWriter(messageType)
		printError(err)
		if _, err := io.Copy(w, r); err != nil {
			fmt.Println(err)
			return
		}
		if err := w.Close(); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func upgradeConnection(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	ws, err := upgrader.Upgrade(w, r, nil)
	printError(err)
	return ws
}

func logMessage(message string) {
	fmt.Println(message)
}

func printError(err error) error {
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
