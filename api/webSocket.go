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

type ClientManager struct {
	mu      sync.Mutex
	clients map[string]*websocket.Conn
}

var clientManager = ClientManager{
	clients: make(map[string]*websocket.Conn),
}

func HandleWebSocket(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username is missing"})
		return
	}

	/*purpose := c.Query("purpose")
	if purpose == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Purpose is missing"})
		return
	}*/

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close()

	clientManager.mu.Lock()
	clientManager.clients[username] = conn
	clientManager.mu.Unlock()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("Received: %s", message)
		if bytes.Equal(message, []byte("exit")) {
			break
		} else if bytes.Equal(message, []byte("update")) {

		}

		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}

	clientManager.mu.Lock()
	delete(clientManager.clients, username)
	clientManager.mu.Unlock()
}

// Function to send a message to a specific client based on username
func sendMessageToClient(username, message string) {
	clientManager.mu.Lock()
	defer clientManager.mu.Unlock()

	conn, exists := clientManager.clients[username]
	if exists {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Write error:", err)
			conn.Close()
			delete(clientManager.clients, username)
		}
	}
}
