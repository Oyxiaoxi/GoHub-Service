package migrations

import (
	"database/sql"
	"GoHub-Service/app/models"
	"GoHub-Service/pkg/migrate"

	"gorm.io/gorm"
)

func init() {

	type User struct {
		models.BaseModel
	}
	type Topic struct {
		models.BaseModel
	}

	type Comment struct {
		models.BaseModel

		TopicID  string `gorm:"type:bigint;not null;index:idx_comments_topic_id"`
		UserID   string `gorm:"type:bigint;not null;index:idx_comments_user_id"`
		Content  string `gorm:"type:text;not null"`
		ParentID string `gorm:"type:bigint;default:0;index:idx_comments_parent_id;comment:父评论ID,0表示顶级评论"`
		LikeCount int64 `gorm:"type:int;default:0;comment:点赞数"`

		// 外键关联
		User  User  `gorm:"foreignKey:UserID"`
		Topic Topic `gorm:"foreignKey:TopicID"`

		models.CommonTimestampsField
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		migrator.AutoMigrate(&Comment{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		migrator.DropTable(&Comment{})
	}

	migrate.Add("2025_12_30_add_comments_table", up, down)
}
