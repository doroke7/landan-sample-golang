package main

import "fmt"

type Man struct {
	Name string
}

type User struct {
	Name string
}

type Car struct {
	Name string
}

func createMan(pMan *Man) {
	*pMan = Man{
		Name: "harry-man",
	}
}

func createUser(pPUser **User) {
	*pPUser = &User{
		Name: "Tom-user",
	}
}

func createCar(oCar *Car) {
	oCar = &Car{
		Name: "bike-car",
	}
}

func main() {

	oMan := Man{
		Name: "TG",
	}

	/*
	   口訣：函數傳值是某個 value 的指標，在函數裡面再指標運算後再改
	*/
	createMan(&oMan)
	fmt.Println(oMan.Name)

	// x1122     x4466    {}
	// **User -> *User -> User

	var pUser *User

	createUser(&pUser)

	fmt.Println(pUser.Name)

	var oCar *Car
	createCar(oCar)

}
