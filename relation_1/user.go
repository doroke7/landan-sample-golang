package main

// User 是關聯範例的主角：一個 User 擁有多筆 UserOrder（has many），
// 外鍵放在 UserOrder 身上（UserOrder.UserId 指回 User.Id）。
type User struct {
	Id         uint
	Name       string
	Age        int
	UserOrders []UserOrder `gorm:"foreignKey:UserId"`
}
