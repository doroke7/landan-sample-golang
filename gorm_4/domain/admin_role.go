package domain

import "time"

// AdminRole 從 internal/domain/admin_role.go 複製過來的最小子集。
type AdminRole struct {
	Id        uint `gorm:"primaryKey"`
	Key       string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
