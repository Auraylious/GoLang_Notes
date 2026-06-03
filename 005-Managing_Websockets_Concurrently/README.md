# Websocket Concurrency
lets set up a functional websocket server, one that can both manage clients, handle events, and broadcast messeges to all the connected clients.        

### Stray, Anonymous, Websockets
by default, when we setup a basic server/client structure with gorilla websockets, there isnt a great way to handle concurrency, broadcast messages, or even identify which socket belongs to which client from other sockets and so on....    

### The Solution
We need to maintain a collection of active connections with identifying metadata that we can iterate or sort through when needed.    
this would allow us not only to broadcast messages to all clients, but also broadcast messages to certain clients as well.    

### The Setup
first we need to setup a go module and install gorilla/websocket in order to use it    

```
go mod init WebsocketTesting
go get github.com/gorilla/websocket
```

we will then build a main.go entry point that decides whether we are running a client or a server      

[main.go](Project/main.go)
```go
//main.go: this file parses our commandline args and gets us where we need to be

package main
import (
    "flag"
    "WebsocketTesting/server"
    "WebsocketTesting/client"
)
func main() {

        // handle commandline arguments
        addressPtr := flag.String("address", "127.0.0.1", "Server to connect to")
        portPtr := flag.String("port", "3000", "Port to use")
        daemonPtr := flag.Bool("daemon", false, "Run as a Server")

        flag.Parse()

        // do the things
        if *daemonPtr {
                // server logic
                server.Listen(*portPtr)
        } else {
                // client logic
                client.Connect(*addressPtr, *portPtr)
        }
}

```

Our Client will be a basic one, it sends Hello World upon connecting and displays any messages received     

[client.go](Project/client/client.go)
```go
package client
import (
	"log"
	"github.com/gorilla/websocket"
)

func Connect(address string, port string) {

	// Connect to Server
	log.Println("connecting to " + address + ":" + port)
        conn, resp, err := websocket.DefaultDialer.Dial("ws://" + address + ":" + port + "/ws", nil)
        if err != nil { panic(err) }

	defer conn.Close()

	// Send Hello World
	log.Println("Sending Hello World")
	err := conn.WriteMessage(websocket.TextMessage, []byte(`Hello World`))
	if err != nil { panic(err) }

	// Process Incoming Messages
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error (connection likely closed): %v", err)
			return
		}
		log.Printf("Recieved: %s\n", message)
	}

}
```

now we need to build our server, we will use this as a starting point:    

```go
package server
import (
    "log"
    "sync"
    "net/http"
    "github.com/gorilla/websocket"
)
var upgrader = websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool { return true },
}
func Listen(port string) {
    router = http.NewServeMux()
    router.HandleFunc("GET /ws", func (w http.ResponseWriter, r *http.Request){
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil { panic(err) }
        defer conn.Close()
        for {
            messageType, data, err := conn.ReadMessage()
            if err != nil { log.Println(err); return }
            log.Printf("Received: %s\n", data)
        }
    })
    log.Fatal(http.ListenAndServe(":" + port, router))
}
```

currently it just accepts connections and prints anything recieved.     

### Building Our Servers Connection Manager

well need to define a connection manager to manage our connections safely with concurrency     
we use a mutex (`sync.Mutex`) to safely lock and unlock the connection manager so things run safely:    

```go
// SocketManager Archetype: Manages active clients safely across multiple go routines
type SocketManager struct {
        sync.Mutex
        clients map[*websocket.Conn]bool
}

// Create a SocketManager object to manage and store our connections
var manager = SocketManager {
        clients: make(map[*websocket.Conn]bool),
}
```

Now when we recieve a new connection we must attach it to the manager like so:      

```go
// Upgrade Request To A WebSocket Connection
conn, err := upgrader.Upgrade(w, r, nil)
if err != nil { fmt.Println(err); return }
defer conn.Close()

// Register the new Socket Safely
manager.Lock()
manager.clients[conn] = true
manager.Unlock()
```

And upon exiting this connection we must remove it from the manager as well:
```go
// Unregister the socket upon disconnection (ending of the previous for loop)
manager.Lock()
delete(manager.clients, conn)
manager.Unlock()
```


Now we need a function to Broadcast our Messages to all the clients:    
```go
// Function to broadcast to all clients
func BroadcastMessage(messageType int, message []byte) {
        // lock the manager until were done broadcasting
        manager.Lock()
        defer manager.Unlock()

        // Iterate over all connected sockets
        for client := range manager.clients {

                // send the message
                err := client.WriteMessage(messageType, message)

                // if we were unable to write to it, close the connection and remove the socket
                if err != nil {
                        log.Printf("Write error to %s, closing connection: %v", client.RemoteAddr(), err)
                        client.Close()
                        delete(manager.clients, client)
                }
        }
}
```

our new server looks like this:    
[server.go](Project/server/server.go)    

```go
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
```


### Putting it All Together

Now when we run our server with `go run main.go --daemon` it binds to port 3000 and listens for connections     
and our client with `go run main.go` it connects to localhost:3000 on the ws address    
we can connect multiple clients at once without having to worry about race conditions because our manager handles them appropriately, instead of things all at once.    

the client then sends `Hello World` to the server, and the server broadcasts it to all connected clients.    

at this point we just need to add metadata to each client connection to identify client1 from client2 and then we can send messages between them as well    
