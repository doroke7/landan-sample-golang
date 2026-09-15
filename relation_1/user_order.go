package main

/*
外鍵口訣：User 有好多 UserOrder，UserOrder 指向 User，UserOrder.UserId 是外鍵。

	User (has many) ──────▶ UserOrder (belongs to)
	  Id                       UserId  ← 外鍵
	  Orders []UserOrder                ← User 端用 foreignKey:UserId 宣告
*/

// UserOrder 屬於某個 User，外鍵是自己身上的 UserId。
type UserOrder struct {
	Id       uint
	UserId   uint
	Product  string
	Quantity int
}
