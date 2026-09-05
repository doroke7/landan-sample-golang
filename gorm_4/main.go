package main

import (
	"context"
	"fmt"
	"log"
	"time"

	bootstrap "example/sample/gorm_4/bootstrap"
	domain "example/sample/gorm_4/domain"
	logic "example/sample/gorm_4/logic"
)

// ============================================================
// 跟主專案一樣按目錄分類，但整套是自包含的，不 import example/bootstrap、
// example/internal/domain、example/internal/output/application/mysql/logic：
//
//	sample/gorm_4/bootstrap  對應主專案的 bootstrap（DB 連線）
//	sample/gorm_4/domain     對應主專案的 internal/domain（AdminUser / AdminRole）
//	sample/gorm_4/logic      對應主專案的 internal/output/application/mysql/logic（AdminUserLogic）
//	sample/gorm_4/main.go    只負責組裝、跑流程
//
// 跑之前要：
//  1. 改 bootstrap/mysql.go 裡的 sMysql* 常數成你本機/測試庫的連線資訊
//  2. 資料庫已經跑過 script/mysql 底下的建表 SQL
//     （tx-admin_users / tx-admin_roles / tx-admin_users_to_admin_roles）
// ============================================================

func main() {
	oDB, oErr := bootstrap.NewMysql()
	if oErr != nil {
		log.Fatal(oErr)
	}

	oContext := context.Background()
	oAdminUserLogic := &logic.AdminUserLogic{DB: oDB, Context: oContext}

	// ===== 1. AdminRole：先查現有角色，等等挑一個綁給新建的 AdminUser =====
	var aAdminRoles []domain.AdminRole
	if oErr := oDB.WithContext(oContext).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		Find(&aAdminRoles).Error; oErr != nil {
		log.Fatal(oErr)
	}

	fmt.Println("[AdminRole] 目前角色：")
	for _, oRole := range aAdminRoles {
		fmt.Printf("  id=%d key=%s name=%s\n", oRole.Id, oRole.Key, oRole.Name)
	}

	var aAdminRoleIds []int
	if len(aAdminRoles) > 0 {
		aAdminRoleIds = append(aAdminRoleIds, int(aAdminRoles[0].Id))
	}

	// ===== 2. AdminUser：透過 logic 層新增一筆，同一個 transaction 裡順便綁角色 =====
	sName := fmt.Sprintf("demo_%d", time.Now().UnixNano())
	sPassword := "12345678"

	oValue := &domain.AdminUserValue{
		Name:         &sName,
		Password:     &sPassword,
		AdminRoleIds: aAdminRoleIds,
	}

	if oErr := oAdminUserLogic.AddAminUser1(oValue); oErr != nil {
		log.Fatal(oErr)
	}
	fmt.Printf("\n[AdminUserLogic.AddAminUser] 新增成功：name=%s\n", sName)
}
