package main

import "fmt"

// 代理模式（Proxy Pattern）
//
// 核心概念：代理與真身實作同一個介面，
// 客戶端透過代理間接存取真身，代理可以在轉發前後插入額外邏輯
// （權限檢查、快取、日誌 ...），對客戶端是透明的。
//
// 本範例：UserService 代理 UserModel。
//   - UserModel   ：真身，真正存取資料
//   - UserService ：代理，轉發給 UserModel，並加上一層存取記錄

// -----------------------------------------------------------------------------
// Subject 介面：代理與真身共用
// -----------------------------------------------------------------------------

type UserRepository interface {
	GetName(iID int) string
}

// -----------------------------------------------------------------------------
// RealSubject：真身
// -----------------------------------------------------------------------------

// UserModel 真正存取資料的物件
type UserModel struct {
	data map[int]string
}

func NewUserModel() *UserModel {
	oModel := &UserModel{
		data: map[int]string{
			1: "Tom",
			2: "Jerry",
		},
	}

	return oModel
}

func (oSelf *UserModel) GetName(iID int) string {
	fmt.Printf("[UserModel] 查詢資料 id=%d\n", iID)

	sName := oSelf.data[iID]

	return sName
}

// -----------------------------------------------------------------------------
// Proxy：與真身實作同一個介面
// -----------------------------------------------------------------------------

// UserService 代理 UserModel
type UserService struct {
	model *UserModel
}

func NewUserService() *UserService {
	oService := &UserService{
		model: NewUserModel(),
	}

	return oService
}

func (oSelf *UserService) GetName(iID int) string {
	fmt.Printf("[UserService] 收到請求 GetName(%d)，轉發給 UserModel\n", iID)

	sName := oSelf.model.GetName(iID)

	fmt.Printf("[UserService] 回傳結果：%q\n", sName)

	return sName
}

// -----------------------------------------------------------------------------

func main() {
	// 客戶端只認得 UserRepository 介面
	var oRepo UserRepository = NewUserService()

	oRepo.GetName(1)

	fmt.Println("---")

	oRepo.GetName(2)
}
