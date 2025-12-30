// Package comment 评论模型
package comment

import (
	"GoHub-Service/app/models"
	"GoHub-Service/app/models/topic"
	"GoHub-Service/app/models/user"
	"GoHub-Service/pkg/database"
)

type Comment struct {
	models.BaseModel

	TopicID   string `json:"topic_id,omitempty"`
	UserID    string `json:"user_id,omitempty"`
	Content   string `json:"content,omitempty"`
	ParentID  string `json:"parent_id,omitempty"`
	LikeCount int64  `json:"like_count,omitempty"`

	// 关联用户
	User user.User `json:"user"`

	// 关联话题
	Topic topic.Topic `json:"topic"`

	models.CommonTimestampsField
}

func (comment *Comment) Create() {
	database.DB.Create(&comment)
}

func (comment *Comment) Save() (rowsAffected int64) {
	result := database.DB.Save(&comment)
	return result.RowsAffected
}

func (comment *Comment) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&comment)
	return result.RowsAffected
}

// GetID 实现Model接口
func (comment *Comment) GetID() uint64 {
	return comment.ID
}

// GetOwnerID 实现OwnershipChecker接口，返回资源所有者ID
func (comment *Comment) GetOwnerID() string {
	return comment.UserID
}
