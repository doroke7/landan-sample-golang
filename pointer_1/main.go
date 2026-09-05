package main

import "fmt"

type AdminUserValue struct {
	Name     *string
	Password *string
}

/*
NOTE: golang 在 struct 運用指標語法
1. . 會隱含多一次 * 解引用運算
2. . 的優先次序比 * 高。
3. 綜合上述 。如 *oAdminUser.Name 其實是 （*(*oAdminUser).Name)
*/
func main() {
	sName := "Tom"
	sPassword := "123456"
	oAdminUser := &AdminUserValue{
		Name:     &sName,
		Password: &sPassword,
	}

	fmt.Println("oAdminUser  是 (AdminUserValue其中一個地址)   =", oAdminUser)
	fmt.Println("*oAdminUser  是 (AdminUserValue其中一個值)   =", *oAdminUser)

	fmt.Println("oAdminUser.Name  是 (～= (*oAdminUser).Name )   =", oAdminUser.Name) // go 的 . 會隱含先 dereference operator（解引用運算子）
	fmt.Println("(*oAdminUser).Name  是 (AdminUserValue值上的 Name 指標)   =", (*oAdminUser).Name)

	fmt.Println("*oAdminUser.Name  是 (～= (*(*oAdminUser).Name) )   =", *oAdminUser.Name) // go 的 . 會隱含先 dereference operator（解引用運算子）

}
