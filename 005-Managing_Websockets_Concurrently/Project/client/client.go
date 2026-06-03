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
