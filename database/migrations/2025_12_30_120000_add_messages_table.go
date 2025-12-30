package migrations

import (
	"database/sql"
	"time"

	"GoHub-Service/app/models"
	"GoHub-Service/pkg/migrate"

	"gorm.io/gorm"
)

func init() {

	type Message struct {
		models.BaseModel

		ConversationID string     `gorm:"type:varchar(64);index"`
		SenderID       string     `gorm:"index;not null"`
		ReceiverID     string     `gorm:"index;not null"`
		Body           string     `gorm:"type:text;not null"`
		ReadAt         *time.Time `gorm:"index"`

		models.CommonTimestampsField
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		migrator.AutoMigrate(&Message{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		migrator.DropTable(&Message{})
	}

	migrate.Add("2025_12_30_120000_add_messages_table", up, down)
}
