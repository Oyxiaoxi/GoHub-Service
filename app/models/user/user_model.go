// Package user 存放用户 Model 相关逻辑
package user

import (
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
