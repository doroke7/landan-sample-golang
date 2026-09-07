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

	data, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
