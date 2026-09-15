package main

import "fmt"

type User struct {
	Id         uint
	Name       string
	Age        int
	UserOrders []UserOrder `gorm:"foreignKey:UserId"`
}

// UserOrder 屬於某個 User，外鍵是自己身上的 UserId。
type UserOrder struct {
	Id       uint
	UserId   uint
	Product  string
	Quantity int
}

// aUserOrders 是原始資料，模擬「資料庫撈出來的一份 flat 訂單」。
var aUserOrders = []UserOrder{
	{Id: 1, UserId: 1, Product: "筆記本", Quantity: 2},
	{Id: 2, UserId: 1, Product: "鉛筆", Quantity: 5},
	{Id: 3, UserId: 1, Product: "橡皮擦", Quantity: 1},
	{Id: 4, UserId: 1, Product: "尺", Quantity: 3},
}

func main() {
	// ===== 1. 手動組裝：直接把 aUserOrders 賦值給 oUser.UserOrders =====
	// 這裡沒有複製，slice 只是複製了 (ptr, len, cap) 這個 header，
	// 底層陣列還是同一個，所以 oUser.UserOrders 跟 aUserOrders 是「別名」關係。
	oUser := User{Id: 1, Name: "Alice", Age: 18}
	oUser.UserOrders = aUserOrders

	fmt.Println("[組裝後，尚未修改]")
	fmt.Println("  oUser.UserOrders[2]:", oUser.UserOrders[2])
	fmt.Println("  aUserOrders[2]:     ", aUserOrders[2])

	// ===== 2. 手動改 oUser.UserOrders 第 3 筆（index 2）=====
	oUser.UserOrders[2].Quantity = 999

	fmt.Println("[改了 oUser.UserOrders[2].Quantity 之後]")
	fmt.Println("  oUser.UserOrders[2]:", oUser.UserOrders[2])
	fmt.Println("  aUserOrders[2]:     ", aUserOrders[2])
	// 因為共用同一個底層陣列，aUserOrders[2] 也被改成 999 了 —— 這是「會一起改掉」。

}
