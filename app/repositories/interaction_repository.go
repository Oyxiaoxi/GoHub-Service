// Package repositories 互动数据访问层
package repositories

import (
	"errors"
	"time"

	"GoHub-Service/app/models/topic"
	"GoHub-Service/app/models/user"
	"GoHub-Service/pkg/database"

	"gorm.io/gorm"
)

// InteractionRepository 定义互动仓储接口
type InteractionRepository interface {
	LikeTopic(userID, topicID string) error
	UnlikeTopic(userID, topicID string) error
	FavoriteTopic(userID, topicID string) error
	UnfavoriteTopic(userID, topicID string) error
	FollowUser(followerID, followeeID string) error
	UnfollowUser(followerID, followeeID string) error
	IncrementTopicView(topicID string) error
}

// TopicLike 话题点赞记录
type TopicLike struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	TopicID   string `gorm:"type:bigint;not null;index:idx_topic_like_topic;uniqueIndex:uidx_topic_like_user_topic"`
	UserID    string `gorm:"type:bigint;not null;index:idx_topic_like_user;uniqueIndex:uidx_topic_like_user_topic"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (TopicLike) TableName() string {
	return "topic_likes"
}

// TopicFavorite 话题收藏记录
type TopicFavorite struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	TopicID   string `gorm:"type:bigint;not null;index:idx_topic_fav_topic;uniqueIndex:uidx_topic_fav_user_topic"`
	UserID    string `gorm:"type:bigint;not null;index:idx_topic_fav_user;uniqueIndex:uidx_topic_fav_user_topic"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (TopicFavorite) TableName() string {
	return "topic_favorites"
}

// UserFollow 用户关注关系
type UserFollow struct {
	ID         uint64 `gorm:"primaryKey;autoIncrement"`
	FollowerID string `gorm:"type:bigint;not null;index:idx_user_follow_follower;uniqueIndex:uidx_user_follow_pair"`
	FolloweeID string `gorm:"type:bigint;not null;index:idx_user_follow_followee;uniqueIndex:uidx_user_follow_pair"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (UserFollow) TableName() string {
	return "user_follows"
}

// interactionRepository 实现
type interactionRepository struct{}

// NewInteractionRepository 创建互动仓储实例
func NewInteractionRepository() InteractionRepository {
	return &interactionRepository{}
}

func (r *interactionRepository) LikeTopic(userID, topicID string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var exists int64
		if err := tx.Model(&TopicLike{}).Where("topic_id = ? AND user_id = ?", topicID, userID).Count(&exists).Error; err != nil {
			return err
		}
		if exists > 0 {
			return nil
		}
		if err := tx.Create(&TopicLike{TopicID: topicID, UserID: userID}).Error; err != nil {
			return err
		}
		if err := tx.Model(&topic.Topic{}).Where("id = ?", topicID).UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error; err != nil {
			return err
		}
		// 增加作者积分
		var ownerID string
		if err := tx.Model(&topic.Topic{}).Select("user_id").Where("id = ?", topicID).Scan(&ownerID).Error; err == nil && ownerID != "" {
			_ = tx.Model(&user.User{}).Where("id = ?", ownerID).UpdateColumn("points", gorm.Expr("points + 1")).Error
		}
		return nil
	})
}

func (r *interactionRepository) UnlikeTopic(userID, topicID string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		res := tx.Where("topic_id = ? AND user_id = ?", topicID, userID).Delete(&TopicLike{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return nil
		}
		if err := tx.Model(&topic.Topic{}).
			Where("id = ? AND like_count > 0", topicID).
			UpdateColumn("like_count", gorm.Expr("like_count - 1")).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *interactionRepository) FavoriteTopic(userID, topicID string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var exists int64
		if err := tx.Model(&TopicFavorite{}).Where("topic_id = ? AND user_id = ?", topicID, userID).Count(&exists).Error; err != nil {
			return err
		}
		if exists > 0 {
			return nil
		}
		if err := tx.Create(&TopicFavorite{TopicID: topicID, UserID: userID}).Error; err != nil {
			return err
		}
		if err := tx.Model(&topic.Topic{}).Where("id = ?", topicID).UpdateColumn("favorite_count", gorm.Expr("favorite_count + 1")).Error; err != nil {
			return err
		}
		var ownerID string
		if err := tx.Model(&topic.Topic{}).Select("user_id").Where("id = ?", topicID).Scan(&ownerID).Error; err == nil && ownerID != "" {
			_ = tx.Model(&user.User{}).Where("id = ?", ownerID).UpdateColumn("points", gorm.Expr("points + 2")).Error
		}
		return nil
	})
}

func (r *interactionRepository) UnfavoriteTopic(userID, topicID string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		res := tx.Where("topic_id = ? AND user_id = ?", topicID, userID).Delete(&TopicFavorite{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return nil
		}
		if err := tx.Model(&topic.Topic{}).
			Where("id = ? AND favorite_count > 0", topicID).
			UpdateColumn("favorite_count", gorm.Expr("favorite_count - 1")).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *interactionRepository) FollowUser(followerID, followeeID string) error {
	if followerID == followeeID {
		return errors.New("cannot follow yourself")
	}
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var exists int64
		if err := tx.Model(&UserFollow{}).Where("follower_id = ? AND followee_id = ?", followerID, followeeID).Count(&exists).Error; err != nil {
			return err
		}
		if exists > 0 {
			return nil
		}
		if err := tx.Create(&UserFollow{FollowerID: followerID, FolloweeID: followeeID}).Error; err != nil {
			return err
		}
		if err := tx.Model(&user.User{}).Where("id = ?", followeeID).UpdateColumn("followers_count", gorm.Expr("followers_count + 1")).Error; err != nil {
			return err
		}
		_ = tx.Model(&user.User{}).Where("id = ?", followeeID).UpdateColumn("points", gorm.Expr("points + 2")).Error
		return nil
	})
}

func (r *interactionRepository) UnfollowUser(followerID, followeeID string) error {
	if followerID == followeeID {
		return errors.New("cannot follow yourself")
	}
	return database.DB.Transaction(func(tx *gorm.DB) error {
		res := tx.Where("follower_id = ? AND followee_id = ?", followerID, followeeID).Delete(&UserFollow{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return nil
		}
		if err := tx.Model(&user.User{}).
			Where("id = ? AND followers_count > 0", followeeID).
			UpdateColumn("followers_count", gorm.Expr("followers_count - 1")).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *interactionRepository) IncrementTopicView(topicID string) error {
	return database.DB.Model(&topic.Topic{}).Where("id = ?", topicID).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}
