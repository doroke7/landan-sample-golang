package main

import (
	"fmt"
)

// user1ToUser2Map 把 user1 形狀的 map 轉成 user2 形狀的 map：
//
//	user1: {id, fisrt_name, last_name, age}
//	user2: {id, name(= first + last), age}
func user1ToUser(oUser1 map[string]any) map[string]any {

	sFirstName, _ := oUser1["fisrt_name"].(string)
	sLastName, _ := oUser1["last_name"].(string)

	oUser2 := map[string]any{
		"id":   oUser1["id"],
		"name": sFirstName + sLastName,
		"age":  oUser1["age"],
	}

	return oUser2
}

func main() {

	oUser1 := map[string]any{
		"id":         1,
		"fisrt_name": "Tom",
		"last_name":  "Stupid",
		"age":        50,
	}

	oUser2 := user1ToUser(oUser1)

	fmt.Println(oUser1)
	fmt.Println(oUser2)
}
