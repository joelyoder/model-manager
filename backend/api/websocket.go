package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"model-manager/backend/database"
	"model-manager/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	// Clients maps ClientID -> thread-safe connection wrapper
	Clients      = make(map[string]*ClientConnection)
	ClientsMutex sync.Mutex
	upgrader     = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all for now
		},
	}
)

type ClientConnection struct {
	Conn *websocket.Conn
	mu   sync.Mutex
}

func (c *ClientConnection) WriteJSON(v interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Conn.WriteJSON(v)
}

type ClientMessage struct {
	Type           string `json:"type"` // "complete", "deleted", "error"
	ModelVersionID uint   `json:"model_version_id"`
	ClientID       string `json:"-"` // Added server-side
}

func HandleWebSocket(c *gin.Context) {
	// Authentication
	secret := c.GetHeader("Authorization")
	if secret != os.Getenv("CLIENT_SECRET") {
		// Try query param for easier testing if needed, or strict header
		if secret == "" {
			secret = c.Query("key")
		}
	}

	if secret != os.Getenv("CLIENT_SECRET") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	clientID := c.GetHeader("X-Client-ID")
	if clientID == "" {
		clientID = c.Query("client_id")
	}
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Client-ID header or client_id query param required"})
		return
	}

	// Upgrade connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade websocket: %v", err)
		return
	}
	defer conn.Close()

	clientConn := &ClientConnection{Conn: conn}

	// Register client
	ClientsMutex.Lock()
	Clients[clientID] = clientConn
	ClientsMutex.Unlock()

	log.Printf("Client connected: %s", clientID)

	defer func() {
		ClientsMutex.Lock()
		delete(Clients, clientID)
		ClientsMutex.Unlock()
		log.Printf("Client disconnected: %s", clientID)
	}()

	// Listen loop
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var msg ClientMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error parsing message from %s: %v", clientID, err)
			continue
		}

		msg.ClientID = clientID
		handleClientMessage(msg)
	}
}

func handleClientMessage(msg ClientMessage) {
	if msg.ModelVersionID == 0 {
		return
	}

	if msg.Type == "complete" {
		// Update/Create ClientFile record
		var cf models.ClientFile
		result := database.DB.Where("client_id = ? AND model_version_id = ?", msg.ClientID, msg.ModelVersionID).First(&cf)
		if result.Error != nil {
			cf = models.ClientFile{
				ClientID:       msg.ClientID,
				ModelVersionID: msg.ModelVersionID,
			}
		}
		cf.Status = "installed"
		database.DB.Save(&cf)
		log.Printf("Updated status 'installed' for model %d on client %s", msg.ModelVersionID, msg.ClientID)

	} else if msg.Type == "deleted" {
		// Remove record or set to something else? Requirement says "Delete Update ClientFile record to remove the entry"
		database.DB.Delete(&models.ClientFile{}, "client_id = ? AND model_version_id = ?", msg.ClientID, msg.ModelVersionID)
		log.Printf("Removed record for model %d on client %s", msg.ModelVersionID, msg.ClientID)
	}
}

// Helper to send to specific client
func SendToClient(clientID string, payload interface{}) error {
	ClientsMutex.Lock()
	clientConn, ok := Clients[clientID]
	ClientsMutex.Unlock()

	if !ok || clientConn == nil {
		return &ClientNotFoundError{ClientID: clientID}
	}

	return clientConn.WriteJSON(payload)
}

type ClientNotFoundError struct {
	ClientID string
}

func (e *ClientNotFoundError) Error() string {
	return "Client not connected: " + e.ClientID
}
