package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	jsonString := `{"id":1,"name":"Tom"}`

	var user map[string]string

	err := json.Unmarshal([]byte(jsonString), &user)
	if err != nil {
		panic(err)
	}

	fmt.Println(user)
	fmt.Println(user["name"])
}
