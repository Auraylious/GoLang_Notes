//server/server.go: our websocket server
package server
import (
	"sync"
	"net/http"
	"github.com/gorilla/socket"
)

// SocketManager Archetype: Manages active clients safely across multiple goroutines
type SocketManager struct {
        sync.Mutex
        clients map[*websocket.Conn]bool
}

// Create a SocketManager to store our connections
var manager = SocketManager {
        clients: make(map[*websocket.Conn]bool),
}

// Function to broadcast to all clients
func BroadcastMessage(messageType int, message []byte) {
	// lock the manager until were done broadcasting
        manager.Lock()
        defer manager.Unlock()

        // Iterate over all connected sockets
        for client := range manager.clients {
                err := client.WriteMessage(messageType, message)
		// if we were unable to write to it, close the connection.
                if err != nil {
			log.Printf("Write error to %s, closing connection: %v", client.RemoteAddr(), err)
                        client.Close()
                        delete(manager.clients, client) // Evict dead sockets instantly
                }
        }
}

// Websocket Upgrader, accepts connections from any origin.
var upgrader = websocket.Upgrader{
        ReadBufferSize: 1024,
        WriteBufferSize: 1024,
        CheckOrigin: func(r *http.Request) bool { return true },
}

// Server Setup
func Listen(port string) {

	router := http.NewServeMux()

        router.HandleFunc("GET /ws", func (w http.ResponseWriter, r *http.Request){

		// Upgrade Request To A WebSocket Connection
                conn, err := upgrader.Upgrade(w, r, nil)
                if err != nil { fmt.Println(err); return }
                defer conn.Close()

		// Register the new Socket Safely
		manager.Lock()
		manager.clients[conn] = true
		manager.Unlock()

		// Notify All Clients that a new client has connected.
		BroadcastMessage("Client connected. Total clients: %d", len(manager.clients))

		// Process Client Messages
		for {
			// recieve message
			messageType, data, err := conn.ReadMessage()
			if err != nil { fmt.Println(err); return }
			fmt.Printf("Received: %s\n", data)

			// send it to all clients
			BroadcastMessage(messageType, data)
			if err != nil { fmt.Println(err); return }

		}

		// Unregister the socket upon disconnection (ending of the previous for loop)
		manager.Lock()
		delete(manager.clients, conn)
		manager.Unlock()
	}
}
