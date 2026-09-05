package bootstrap

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 從 bootstrap/mysql.go 複製過來的最小版本，連線資訊改成下面幾個常數，自己填。
const (
	sMysqlUser   = "root"
	sMysqlPass   = "root"
	sMysqlHost   = "127.0.0.1"
	sMysqlPort   = "3306"
	sMysqlDbName = "example"
	sTablePrefix = "tx-"
)

func NewMysql() (*gorm.DB, error) {
	sDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		sMysqlUser, sMysqlPass, sMysqlHost, sMysqlPort, sMysqlDbName,
	)

	return gorm.Open(mysql.Open(sDSN), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: sTablePrefix,
		},
	})
}
