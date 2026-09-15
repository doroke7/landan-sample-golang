package main

import "fmt"

// 外鍵放在 UserOrder 身上（UserOrder.UserId 指回 User.Id）。
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

func main() {

	oUser := User{
		Id:   1,
		Name: "Alice",
		Age:  18,
		UserOrders: []UserOrder{
			{Id: 1, UserId: 1, Product: "筆記本", Quantity: 2},
			{Id: 2, UserId: 1, Product: "鉛筆", Quantity: 5},
			{Id: 3, UserId: 1, Product: "橡皮擦", Quantity: 1},
			{Id: 4, UserId: 1, Product: "尺", Quantity: 3},
		},
	}

	var oUserOrder = oUser.UserOrders[1]

	oUserOrder.Quantity = 100

	fmt.Println("oUserOrder=", oUserOrder)
	fmt.Println("oUser.UserOrders[1]=", oUser.UserOrders[1])

	// 關係裡面的 slice 元素 如果是 值 => 不會互相修改
	// 關係裡面的 slice 框 如果是 值   => 不會互相修改

}
