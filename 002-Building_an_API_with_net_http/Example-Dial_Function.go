package main;

import (
	"fmt"
	"net"
	"bufio"
)

func main(){
	//Example HTTP Get / request via Dial Function

	// open a tcp connection
	conn, err := net.Dial("tcp", "golang.org:80")

	if err != nil {
	    // handle any errors from opening the connection
	    fmt.Println("There was an Error")
	    fmt.Println(err)
	}

	// send a get request
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")

	// read the result from the buffer
	status, err := bufio.NewReader(conn).ReadString('\n')

	// print to terminal
	fmt.Println("Status: " + status)
}
