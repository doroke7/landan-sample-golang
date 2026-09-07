package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	jsonString := `{"id":1,"name":"Tom"}`

	var user User

	err := json.Unmarshal([]byte(jsonString), &user)
	if err != nil {
		panic(err)
	}

	fmt.Println(user)
	fmt.Println(user.Name)
}
