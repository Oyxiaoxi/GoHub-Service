package migrations

import (
	"GoHub-Service/app/models"
	"GoHub-Service/pkg/migrate"
	"database/sql"

	"gorm.io/gorm"
)

func init() {
	type Topic struct {
		models.BaseModel
		LikeCount     int64 `gorm:"type:int;default:0;index:idx_topics_like_count"`
		FavoriteCount int64 `gorm:"type:int;default:0;index:idx_topics_favorite_count"`
		ViewCount     int64 `gorm:"type:int;default:0;index:idx_topics_view_count"`
	}

	type User struct {
		models.BaseModel
		FollowersCount int64 `gorm:"type:int;default:0;index:idx_users_followers_count"`
		Points         int64 `gorm:"type:int;default:0;index:idx_users_points"`
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.AutoMigrate(&Topic{}, &User{})
		// AddColumn 忽略已存在的列，AutoMigrate 已处理字段
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.DropColumn(&Topic{}, "LikeCount")
		_ = migrator.DropColumn(&Topic{}, "FavoriteCount")
		_ = migrator.DropColumn(&Topic{}, "ViewCount")
		_ = migrator.DropColumn(&User{}, "FollowersCount")
		_ = migrator.DropColumn(&User{}, "Points")
	}

	migrate.Add("2025_12_30_010000_add_interaction_columns", up, down)
}
