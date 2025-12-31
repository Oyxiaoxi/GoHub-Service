package follow

import (
	"GoHub-Service/app/models"
	"GoHub-Service/pkg/database"
)

// Follow 关注关系模型
type Follow struct {
	models.BaseModel
	UserID   string `gorm:"index:idx_user_follow;not null" json:"user_id"`        // 关注者ID
	FollowID string `gorm:"index:idx_user_follow;not null" json:"follow_id"`      // 被关注者ID
	models.CommonTimestampsField
}

// TableName 指定表名
func (Follow) TableName() string {
	return "follows"
}

// Create 创建关注
func (f *Follow) Create() {
	database.DB.Create(&f)
}

// Delete 删除关注
func (f *Follow) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&f)
	return result.RowsAffected
}

// GetByUserID 获取用户关注的人
func GetByUserID(userID string) []Follow {
	var follows []Follow
	database.DB.Where("user_id = ?", userID).Find(&follows)
	return follows
}

// GetFollowersCount 获取粉丝数
func GetFollowersCount(userID string) int64 {
	var count int64
	database.DB.Model(&Follow{}).Where("follow_id = ?", userID).Count(&count)
	return count
}

// GetFollowingCount 获取关注数
func GetFollowingCount(userID string) int64 {
	var count int64
	database.DB.Model(&Follow{}).Where("user_id = ?", userID).Count(&count)
	return count
}

// HasFollowed 检查是否已关注
func HasFollowed(userID, followID string) bool {
	var count int64
	database.DB.Model(&Follow{}).
		Where("user_id = ? AND follow_id = ?", userID, followID).
		Count(&count)
	return count > 0
}
