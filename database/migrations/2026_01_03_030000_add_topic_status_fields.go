package migrations

import (
	"database/sql"

	"GoHub-Service/pkg/migrate"

	"gorm.io/gorm"
)

func init() {
	type Topic struct {
		Status       int    `gorm:"type:int;default:1;index;comment:状态:0待审核,1已通过,-1已拒绝"`
		RejectReason string `gorm:"type:varchar(500);comment:拒绝原因"`
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.AutoMigrate(&Topic{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.DropColumn(&Topic{}, "status")
		_ = migrator.DropColumn(&Topic{}, "reject_reason")
	}

	migrate.Add("2026_01_03_030000_add_topic_status_fields", up, down)
}
