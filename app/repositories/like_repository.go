package repositories

import (
	"GoHub-Service/app/models/like"
	"GoHub-Service/pkg/database"
)

// LikeRepository 点赞仓储接口
type LikeRepository interface {
	// Create 创建点赞
	Create(l *like.Like) error
	// Delete 删除点赞
	Delete(userID, targetType, targetID string) error
	// GetByID 根据ID获取点赞
	GetByID(id string) (*like.Like, error)
	// GetByTarget 获取目标的点赞列表
	GetByTarget(targetType, targetID string, offset, limit int) ([]like.Like, int64, error)
	// GetAll 获取所有点赞
	GetAll(offset, limit int) ([]like.Like, int64, error)
	// Count 获取点赞总数
	Count() (int64, error)
	// CountByTarget 获取目标的点赞总数
	CountByTarget(targetType, targetID string) (int64, error)
	// CountByUser 获取用户的点赞总数
	CountByUser(userID string) (int64, error)
	// HasLiked 检查是否已点赞
	HasLiked(userID, targetType, targetID string) (bool, error)
}

type LikeRepositoryImpl struct{}

// NewLikeRepository 创建点赞仓储
func NewLikeRepository() LikeRepository {
	return &LikeRepositoryImpl{}
}

// Create 创建点赞
func (r *LikeRepositoryImpl) Create(l *like.Like) error {
	if err := database.DB.Create(l).Error; err != nil {
		return err
	}
	return nil
}

// Delete 删除点赞
func (r *LikeRepositoryImpl) Delete(userID, targetType, targetID string) error {
	if err := database.DB.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).
		Delete(&like.Like{}).Error; err != nil {
		return err
	}
	return nil
}

// GetByID 根据ID获取点赞
func (r *LikeRepositoryImpl) GetByID(id string) (*like.Like, error) {
	var l like.Like
	if err := database.DB.Where("id = ?", id).First(&l).Error; err != nil {
		return nil, err
	}
	return &l, nil
}

// GetByTarget 获取目标的点赞列表
func (r *LikeRepositoryImpl) GetByTarget(targetType, targetID string, offset, limit int) ([]like.Like, int64, error) {
	var likes []like.Like
	var count int64

	if err := database.DB.Where("target_type = ? AND target_id = ?", targetType, targetID).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := database.DB.
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&likes).Error; err != nil {
		return nil, 0, err
	}

	return likes, count, nil
}

// GetAll 获取所有点赞
func (r *LikeRepositoryImpl) GetAll(offset, limit int) ([]like.Like, int64, error) {
	var likes []like.Like
	var count int64

	if err := database.DB.Model(&like.Like{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := database.DB.
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&likes).Error; err != nil {
		return nil, 0, err
	}

	return likes, count, nil
}

// Count 获取点赞总数
func (r *LikeRepositoryImpl) Count() (int64, error) {
	var count int64
	if err := database.DB.Model(&like.Like{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByTarget 获取目标的点赞总数
func (r *LikeRepositoryImpl) CountByTarget(targetType, targetID string) (int64, error) {
	var count int64
	if err := database.DB.Model(&like.Like{}).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByUser 获取用户的点赞总数
func (r *LikeRepositoryImpl) CountByUser(userID string) (int64, error) {
	var count int64
	if err := database.DB.Model(&like.Like{}).
		Where("user_id = ?", userID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// HasLiked 检查是否已点赞
func (r *LikeRepositoryImpl) HasLiked(userID, targetType, targetID string) (bool, error) {
	var count int64
	if err := database.DB.Model(&like.Like{}).
		Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
