package domain

import "time"

// AdminUser 從 internal/domain/admin_user.go 複製過來的最小子集。
type AdminUser struct {
	Id         uint `gorm:"primaryKey"`
	Name       string
	Password   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
	AdminRoles []AdminRole `gorm:"many2many:admin_users_to_admin_roles"`
}

type AdminUserValue struct {
	Name         *string
	Password     *string
	AdminRoleIds []int
}

// AdminUsersToAdminRole 純樞紐表：只有兩個外鍵、複合主鍵，沒有 id / 時間欄。
type AdminUsersToAdminRole struct {
	AdminUserId uint `gorm:"primaryKey"`
	AdminRoleId uint `gorm:"primaryKey"`
}
