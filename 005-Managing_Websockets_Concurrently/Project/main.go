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
