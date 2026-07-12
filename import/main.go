package main

import (
	"fmt"
	importModel "landan-backend-grpc/sample/import/model"
)

/*
*

		Golang 的思維是 引入 檔案的全局 變量
	    所以其他的文件 用 var 在 main 外面 宣告一個 大寫變量，就可以被引入
*/
func main() {
	var oTestModel = importModel.TestModel
	fmt.Println(oTestModel)
}
