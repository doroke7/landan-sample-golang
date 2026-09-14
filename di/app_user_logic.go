package main

type AppUserLogic struct {
	*AppUserModel
}

// 錯誤的 DI 範例
func NewAppUserLogic() *AppUserLogic {
	oAppUserModel := NewAppUserModel()

	return &AppUserLogic{
		AppUserModel: oAppUserModel,
	}
}
