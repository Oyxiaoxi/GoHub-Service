// Package user 存放用户 Model 相关逻辑
package user

import (
	"time"

	"GoHub-Service/app/models"
	"GoHub-Service/pkg/database"
	"GoHub-Service/pkg/hash"
)

// User 用户模型
type User struct {
	models.BaseModel

	Name string `json:"name,omitempty"`

	City           string `json:"city,omitempty"`
	Introduction   string `json:"introduction,omitempty"`
	Avatar         string `json:"avatar,omitempty"`
	FollowersCount int64  `json:"followers_count,omitempty"`
	Points         int64  `json:"points,omitempty"`

	Email    string `gorm:"uniqueIndex" json:"-"`
	Phone    string `gorm:"uniqueIndex" json:"-"`
	Password string `json:"-"`

	// 封禁相关字段
	IsBanned  bool       `gorm:"type:boolean;default:false;index;comment:是否封禁" json:"is_banned,omitempty"`
	BannedAt  *time.Time `gorm:"comment:封禁时间" json:"banned_at,omitempty"`
	BannedBy  uint64     `gorm:"comment:封禁操作员ID" json:"banned_by,omitempty"`
	BanReason string     `gorm:"type:varchar(500);comment:封禁原因" json:"ban_reason,omitempty"`
	BanUntil  *time.Time `gorm:"comment:封禁截止时间" json:"ban_until,omitempty"`

	models.CommonTimestampsField
}

// Create 创建用户，通过 User.ID 来判断是否创建成功
func (userModel *User) Create() {
	database.DB.Create(&userModel)
}

// ComparePassword 密码是否正确
func (userModel *User) ComparePassword(_password string) bool {
	return hash.BcryptCheck(_password, userModel.Password)
}

// Save 保存用户实例
func (userModel *User) Save() (rowsAffected int64) {
	result := database.DB.Save(&userModel)
	return result.RowsAffected
}
