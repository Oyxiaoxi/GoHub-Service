package migrations

import (
	"GoHub-Service/app/models"
	"GoHub-Service/pkg/migrate"
	"database/sql"

	"gorm.io/gorm"
)

func init() {
	type TopicLike struct {
		models.BaseModel
		TopicID string `gorm:"type:bigint;not null;index:idx_topic_like_topic"`
		UserID  string `gorm:"type:bigint;not null;index:idx_topic_like_user;uniqueIndex:uidx_topic_like_user_topic"`
		// unique on user-topic
		_ struct{} `gorm:"uniqueIndex:uidx_topic_like_user_topic"`
		models.CommonTimestampsField
	}

	type TopicFavorite struct {
		models.BaseModel
		TopicID string   `gorm:"type:bigint;not null;index:idx_topic_fav_topic"`
		UserID  string   `gorm:"type:bigint;not null;index:idx_topic_fav_user;uniqueIndex:uidx_topic_fav_user_topic"`
		_       struct{} `gorm:"uniqueIndex:uidx_topic_fav_user_topic"`
		models.CommonTimestampsField
	}

	type UserFollow struct {
		models.BaseModel
		FollowerID string   `gorm:"type:bigint;not null;index:idx_user_follow_follower;uniqueIndex:uidx_user_follow_pair"`
		FolloweeID string   `gorm:"type:bigint;not null;index:idx_user_follow_followee;uniqueIndex:uidx_user_follow_pair"`
		_          struct{} `gorm:"uniqueIndex:uidx_user_follow_pair"`
		models.CommonTimestampsField
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.AutoMigrate(&TopicLike{}, &TopicFavorite{}, &UserFollow{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.DropTable(&TopicLike{})
		_ = migrator.DropTable(&TopicFavorite{})
		_ = migrator.DropTable(&UserFollow{})
	}

	migrate.Add("2025_12_30_010100_create_interaction_tables", up, down)
}
