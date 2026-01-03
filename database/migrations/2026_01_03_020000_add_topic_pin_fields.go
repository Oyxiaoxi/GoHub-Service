package migrations

import (
	"database/sql"
	"time"

	"GoHub-Service/pkg/migrate"

	"gorm.io/gorm"
)

func init() {
	type Topic struct {
		IsPinned bool       `gorm:"type:boolean;default:false;index;comment:是否置顶"`
		PinnedAt *time.Time `gorm:"comment:置顶时间"`
		PinnedBy uint64     `gorm:"comment:置顶操作员ID"`
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.AutoMigrate(&Topic{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.DropColumn(&Topic{}, "is_pinned")
		_ = migrator.DropColumn(&Topic{}, "pinned_at")
		_ = migrator.DropColumn(&Topic{}, "pinned_by")
	}

	migrate.Add("2026_01_03_020000_add_topic_pin_fields", up, down)
}
