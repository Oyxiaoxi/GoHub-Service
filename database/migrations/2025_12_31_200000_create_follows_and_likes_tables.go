package migrations

import (
	"GoHub-Service/app/models"
	"GoHub-Service/pkg/migrate"
	"database/sql"

	"gorm.io/gorm"
)

func init() {
	// Follow 关注关系表
	type Follow struct {
		models.BaseModel
		UserID   string `gorm:"type:varchar(255);not null;index:idx_user_follow;comment:关注者ID"`
		FollowID string `gorm:"type:varchar(255);not null;index:idx_user_follow;comment:被关注者ID"`
		models.CommonTimestampsField
	}

	// Like 点赞表
	type Like struct {
		models.BaseModel
		UserID     string `gorm:"type:varchar(255);not null;index:idx_like_user;comment:点赞用户ID"`
		TargetType string `gorm:"type:varchar(50);not null;index:idx_like_target;comment:目标类型(topic/comment)"`
		TargetID   string `gorm:"type:varchar(255);not null;index:idx_like_target;comment:目标ID"`
		models.CommonTimestampsField
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		// 创建 follows 表
		_ = migrator.AutoMigrate(&Follow{})
		
		// 创建 likes 表
		_ = migrator.AutoMigrate(&Like{})
		
		// 添加复合唯一索引
		_, _ = DB.Exec("ALTER TABLE follows ADD UNIQUE INDEX uidx_follow_pair (user_id, follow_id)")
		_, _ = DB.Exec("ALTER TABLE likes ADD UNIQUE INDEX uidx_like_target_user (user_id, target_type, target_id)")
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.DropTable("follows")
		_ = migrator.DropTable("likes")
	}

	migrate.Add("2025_12_31_200000_create_follows_and_likes_tables", up, down)
}
