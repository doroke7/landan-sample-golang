package main

import (
	"fmt"
)

type User1 struct {
	Id        int    `json:"id"`
	FisrtName string `json:"fisrt_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

type User2 struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func user1ToUser2(oUser1 User1) User2 {

	oUser2 := User2{
		Id:   oUser1.Id,
		Name: oUser1.FisrtName + oUser1.LastName,
		Age:  oUser1.Age,
	}
	return oUser2
}

func main() {

	oUser1 := User1{
		Id:        1,
		FisrtName: "Tom",
		LastName:  "Stupid",
		Age:       50,
	}

	oUser2 := user1ToUser2(oUser1)

	fmt.Println(oUser2)
}
