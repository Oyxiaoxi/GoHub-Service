// Package topic 模型
package topic

import (
	"time"

	"GoHub-Service/app/models"
	"GoHub-Service/app/models/category"
	"GoHub-Service/app/models/user"
	"GoHub-Service/pkg/database"

	"github.com/spf13/cast"
)

type Topic struct {
	models.BaseModel

	Title         string `json:"title,omitempty"`
	Body          string `json:"body,omitempty"`
	UserID        string `gorm:"index" json:"user_id,omitempty"`
	CategoryID    string `gorm:"index" json:"category_id,omitempty"`
	LikeCount     int64  `json:"like_count,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	ViewCount     int64  `json:"view_count,omitempty"`

	// 置顶相关字段
	IsPinned bool       `gorm:"type:boolean;default:false;index;comment:是否置顶" json:"is_pinned,omitempty"`
	PinnedAt *time.Time `gorm:"comment:置顶时间" json:"pinned_at,omitempty"`
	PinnedBy uint64     `gorm:"comment:置顶操作员ID" json:"pinned_by,omitempty"`

	// 审核相关字段
	Status       int    `gorm:"type:int;default:1;index;comment:状态:0待审核,1已通过,-1已拒绝" json:"status,omitempty"`
	RejectReason string `gorm:"type:varchar(500);comment:拒绝原因" json:"reject_reason,omitempty"`

	// 通过 user_id 关联用户
	User user.User `json:"user"`

	// 通过 category_id 关联分类
	Category category.Category `json:"category"`

	models.CommonTimestampsField
}

func (topic *Topic) Create() {
	database.DB.Create(&topic)
}

func (topic *Topic) Save() (rowsAffected int64) {
	result := database.DB.Save(&topic)
	return result.RowsAffected
}

func (topic *Topic) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&topic)
	return result.RowsAffected
}

// GetID 实现Model接口
func (topic *Topic) GetID() uint64 {
	return topic.ID
}

// GetOwnerID 实现OwnershipChecker接口，返回资源所有者ID
func (topic *Topic) GetOwnerID() string {
	return topic.UserID
}

// GetStringID 获取字符串格式的ID
func (topic *Topic) GetStringID() string {
	return cast.ToString(topic.ID)
}
