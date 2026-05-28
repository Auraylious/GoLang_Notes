package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func main() {
	// Marshalling: Go Struct -> JSON
	u := User{ID: 1, Username: "gopher"}
	jsonData, _ := json.Marshal(u)
	fmt.Println(string(jsonData)) // {"id":1,"username":"gopher"}

	// Unmarshalling: JSON -> Go Struct
	input := `{"id":2,"username":"newuser"}`
	var decodedUser User
	json.Unmarshal([]byte(input), &decodedUser)
	fmt.Printf("%+v\n", decodedUser) // {ID:2 Username:newuser}
}
