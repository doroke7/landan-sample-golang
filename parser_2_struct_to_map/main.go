package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	user := User{
		ID:   1,
		Name: "Tom",
		Age:  18,
	}

	// struct -> map：先 Marshal 成 JSON bytes，再 Unmarshal 進 map（會照 json tag 命名）
	data, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	var userMap map[string]any

	err = json.Unmarshal(data, &userMap)
	if err != nil {
		panic(err)
	}

	fmt.Println(userMap)
	fmt.Println(userMap["name"])
}
