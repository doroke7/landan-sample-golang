package main

/*
   外鍵口訣1. User 有好多 Post, post 指向 user， post.user_id 是外鍵
   外鍵口訣2. User 屬於一個 Nation, User 指向 Nation， user.nation_id 是外鍵


	外鍵口訣3. has one 是 has many 的特例

*/

type User struct {
	Id       uint
	NationId uint
	Name     string
	Age      int

	// 1. belongs to：User 屬於某個 Nation，外鍵是自己身上的 NationId
	Nation Nation `gorm:"foreignKey:NationId"`

	// 2. has one：一個 User 有一份 UserDetail，外鍵是 UserDetail.UserId
	Detail UserDetail `gorm:"foreignKey:UserId"`

	// 3. has many
	Posts []Post `gorm:"foreignKey:UserId"`
}

type UserDetail struct {
	Id          uint
	UserId      uint
	Description string
}

type Post struct {
	Id     uint
	UserId uint
	Text   string
}

type Nation struct {
	Id          uint
	Name        string
	Description string

	// 3. has many：一個 Nation 有多個 User，外鍵是 User.NationId
	Users []User `gorm:"foreignKey:NationId"`
}

func main() {

	//
}
