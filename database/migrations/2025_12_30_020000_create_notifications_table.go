package migrations

import (
	"GoHub-Service/app/models"
	"GoHub-Service/pkg/migrate"
	"database/sql"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func init() {
	type Notification struct {
		models.BaseModel
		UserID  string         `gorm:"type:bigint;not null;index:idx_notifications_user"`
		ActorID string         `gorm:"type:bigint;index:idx_notifications_actor"`
		Type    string         `gorm:"type:varchar(50);not null;index:idx_notifications_type"`
		Data    datatypes.JSON `gorm:"type:json"`
		ReadAt  *sql.NullTime  `gorm:"index:idx_notifications_read_at"`
		models.CommonTimestampsField
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.AutoMigrate(&Notification{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.DropTable(&Notification{})
	}

	migrate.Add("2025_12_30_020000_create_notifications_table", up, down)
}
