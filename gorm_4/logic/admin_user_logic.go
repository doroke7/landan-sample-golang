package logic

import (
	"context"
	"errors"

	"gorm.io/gorm"

	domain "example/sample/gorm_4/domain"
)

// AdminUserLogic 從 internal/output/application/mysql/logic/admin_user_logic.go 複製過來的核心方法。
type AdminUserLogic struct {
	DB      *gorm.DB
	Context context.Context
}

func (oSelf *AdminUserLogic) AddAminUser1(oValue *domain.AdminUserValue) error {
	oColumns := map[string]any{}
	if oValue.Name != nil {
		oColumns["name"] = *oValue.Name
	}
	if oValue.Password != nil {
		oColumns["password"] = *oValue.Password
	}

	oTx := oSelf.DB.WithContext(oSelf.Context).Begin()
	if oTx.Error != nil {
		return oTx.Error
	}

	defer func() {
		if oRecover := recover(); oRecover != nil {
			oTx.Rollback()
		}
	}()

	oResult := oTx.Model(&domain.AdminUser{}).Create(oColumns)
	if oResult.Error != nil {
		oTx.Rollback()
		return oResult.Error
	}
	if oResult.RowsAffected == 0 {
		oTx.Rollback()
		oErr := errors.New("新增0筆")

		return oErr
	}

	var iAdminUserId uint64
	oResult = oTx.Raw("SELECT LAST_INSERT_ID()").Scan(&iAdminUserId)
	if oResult.Error != nil {
		oTx.Rollback()
		return oResult.Error
	}

	oResult = oTx.Where("admin_user_id = ?", iAdminUserId).Delete(&domain.AdminUsersToAdminRole{})
	if oResult.Error != nil {
		oTx.Rollback()
		return oResult.Error
	}

	if len(oValue.AdminRoleIds) > 0 {
		aRelations := make([]domain.AdminUsersToAdminRole, 0, len(oValue.AdminRoleIds))
		for _, iAdminRoleId := range oValue.AdminRoleIds {
			aRelations = append(aRelations, domain.AdminUsersToAdminRole{
				AdminUserId: uint(iAdminUserId),
				AdminRoleId: uint(iAdminRoleId),
			})
		}

		oResult = oTx.Create(&aRelations)
		if oResult.Error != nil {
			oTx.Rollback()
			return oResult.Error
		}
	}

	oResult = oTx.Commit()

	return oResult.Error
}

// AddAminUser2 跟 AddAminUser1 做一樣的事，但改用 gorm 的 Transaction(func(tx *gorm.DB) error {...})
// callback 寫法：不用自己手動 Begin/Commit，只要 closure 回傳 error 就自動 Rollback，回傳 nil 才自動 Commit。
func (oSelf *AdminUserLogic) AddAminUser2(oValue *domain.AdminUserValue) error {
	oColumns := map[string]any{}
	if oValue.Name != nil {
		oColumns["name"] = *oValue.Name
	}
	if oValue.Password != nil {
		oColumns["password"] = *oValue.Password
	}

	return oSelf.DB.WithContext(oSelf.Context).Transaction(func(oTx *gorm.DB) error {

		oResult := oTx.Model(&domain.AdminUser{}).Create(oColumns)
		if oResult.Error != nil {
			return oResult.Error
		}
		if oResult.RowsAffected == 0 {
			return errors.New("新增0筆")
		}

		var iAdminUserId uint64
		if oErr := oTx.Raw("SELECT LAST_INSERT_ID()").Scan(&iAdminUserId).Error; oErr != nil {
			return oErr
		}

		if oErr := oTx.Where("admin_user_id = ?", iAdminUserId).Delete(&domain.AdminUsersToAdminRole{}).Error; oErr != nil {
			return oErr
		}

		if len(oValue.AdminRoleIds) == 0 {
			return nil
		}

		aRelations := make([]domain.AdminUsersToAdminRole, 0, len(oValue.AdminRoleIds))
		for _, iAdminRoleId := range oValue.AdminRoleIds {
			aRelations = append(aRelations, domain.AdminUsersToAdminRole{
				AdminUserId: uint(iAdminUserId),
				AdminRoleId: uint(iAdminRoleId),
			})
		}

		return oTx.Create(&aRelations).Error
	})
}
