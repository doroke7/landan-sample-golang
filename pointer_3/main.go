package main

import "fmt"

func main() {

	iA := 5
	// pA := &iA
	test1(&iA)

	iB := 2
	pB := &iB

	test2(&pB)

	fmt.Println(iA)
	fmt.Println(*pB)

}

// 口訣：輸入 *int ，修改  *int 指向的 int ，達到改變
func test1(pNumber *int) {

	*pNumber = 30
}

// 口訣：輸入 **int ，修改  **int 指向的 *int ，達到改變
func test2(pPNumber **int) {

	iInt := 10
	*pPNumber = &iInt
}
