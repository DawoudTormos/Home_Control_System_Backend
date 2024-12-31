package api

import (
	"bytes"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ClientManagerDeviceLinking struct {
	mu      sync.Mutex
	clients map[string]DeviceLinking
}

type ClientManagerDataUpdate struct {
	mu      sync.Mutex
	clients map[string]DataUpdate
}

type DeviceLinking struct {
	deviceID int32
	conn     *websocket.Conn
}

type DataUpdate struct {
	conn *websocket.Conn
}

var clientManagerDeviceLinking = ClientManagerDeviceLinking{
	clients: make(map[string]DeviceLinking),
}

var clientManagerDataUpdate = ClientManagerDataUpdate{
	clients: make(map[string]DataUpdate),
}

func HandleWebSocket(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username is missing"})
		return
	}

	purpose := c.Query("purpose")
	if purpose == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Purpose is missing."})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close()

	if purpose == "device_linking" {

		var d DeviceLinking
		d.conn = conn
		d.deviceID = 0
		clientManagerDeviceLinking.mu.Lock()
		clientManagerDeviceLinking.clients[username] = d
		clientManagerDeviceLinking.mu.Unlock()

		for {

			clientManagerDeviceLinking.mu.Lock()
			_, message, err := clientManagerDeviceLinking.clients[username].conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}
			log.Printf("Received: %s", message)

			if bytes.Equal(message, []byte("exit")) {
				break
			} else if bytes.Equal(message, []byte("update")) {

			}

			err = clientManagerDeviceLinking.clients[username].conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("Write error:", err)
				break
			}
			clientManagerDeviceLinking.mu.Unlock()
		}

		clientManagerDeviceLinking.mu.Lock()
		delete(clientManagerDeviceLinking.clients, username)
		clientManagerDeviceLinking.mu.Unlock()

	} else if purpose == "data_update" {

		var d DataUpdate
		d.conn = conn
		clientManagerDataUpdate.mu.Lock()
		clientManagerDataUpdate.clients[username] = d
		clientManagerDataUpdate.mu.Unlock()

		for {
			_, message, err := clientManagerDataUpdate.clients[username].conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}
			log.Printf("Received: %s", message)

			if bytes.Equal(message, []byte("exit")) {
				break
			} else if bytes.Equal(message, []byte("update")) {

			}

			err = clientManagerDataUpdate.clients[username].conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("Write error:", err)
				break
			}
		}

		clientManagerDataUpdate.mu.Lock()
		delete(clientManagerDataUpdate.clients, username)
		clientManagerDataUpdate.mu.Unlock()

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid purpose"})
		return
	}

}

// Function to send a message to a specific client based on username
func sendMessageToClient(username, message string) {
	clientManagerDataUpdate.mu.Lock()
	defer clientManagerDataUpdate.mu.Unlock()

	conn := clientManagerDataUpdate.clients[username].conn
	if conn != nil {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Write error:", err)
			conn.Close()
			delete(clientManagerDataUpdate.clients, username)
		}
	}
}
