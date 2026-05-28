package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	// 1. Send the GET request
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		fmt.Printf("Request failed: %s\n", err)
		return
	}

	// 2. Ensure the response body is closed to prevent resource leaks
	defer resp.Body.Close()

	// 3. Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read body: %s\n", err)
		return
	}

	// 4. Output the results
	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Body:", string(body))
}
