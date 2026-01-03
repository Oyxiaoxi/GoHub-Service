package migrations

import (
	"database/sql"
	"time"

	"GoHub-Service/pkg/migrate"

	"gorm.io/gorm"
)

func init() {
	type User struct {
		IsBanned  bool       `gorm:"type:boolean;default:false;index;comment:是否封禁"`
		BannedAt  *time.Time `gorm:"comment:封禁时间"`
		BannedBy  uint64     `gorm:"comment:封禁操作员ID"`
		BanReason string     `gorm:"type:varchar(500);comment:封禁原因"`
		BanUntil  *time.Time `gorm:"comment:封禁截止时间"`
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.AutoMigrate(&User{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.DropColumn(&User{}, "is_banned")
		_ = migrator.DropColumn(&User{}, "banned_at")
		_ = migrator.DropColumn(&User{}, "banned_by")
		_ = migrator.DropColumn(&User{}, "ban_reason")
		_ = migrator.DropColumn(&User{}, "ban_until")
	}

	migrate.Add("2026_01_03_010000_add_user_ban_fields", up, down)
}
