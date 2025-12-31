package like

import (
	"GoHub-Service/app/models"
	"GoHub-Service/pkg/database"
)

// Like 点赞模型
type Like struct {
	models.BaseModel
	UserID     string `gorm:"index:idx_like_target;not null" json:"user_id"`      // 点赞者ID
	TargetType string `gorm:"index:idx_like_target;not null" json:"target_type"`  // 目标类型 (topic/comment)
	TargetID   string `gorm:"index:idx_like_target;not null" json:"target_id"`    // 目标ID
	models.CommonTimestampsField
}

// TableName 指定表名
func (Like) TableName() string {
	return "likes"
}

// Create 创建点赞
func (l *Like) Create() {
	database.DB.Create(&l)
}

// Delete 删除点赞
func (l *Like) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&l)
	return result.RowsAffected
}

// GetCount 获取点赞数
func GetCount(targetType, targetID string) int64 {
	var count int64
	database.DB.Model(&Like{}).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Count(&count)
	return count
}

// HasLiked 检查是否已点赞
func HasLiked(userID, targetType, targetID string) bool {
	var count int64
	database.DB.Model(&Like{}).
		Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).
		Count(&count)
	return count > 0
}
