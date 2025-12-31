package repositories

import (
	"GoHub-Service/app/models/follow"
	"GoHub-Service/pkg/database"
)

// FollowRepository 关注仓储接口
type FollowRepository interface {
	// Create 创建关注关系
	Create(f *follow.Follow) error
	// Delete 删除关注关系
	Delete(userID, followID string) error
	// GetByID 根据ID获取关注
	GetByID(id string) (*follow.Follow, error)
	// GetFollowers 获取粉丝列表
	GetFollowers(userID string, offset, limit int) ([]follow.Follow, int64, error)
	// GetFollowing 获取关注列表
	GetFollowing(userID string, offset, limit int) ([]follow.Follow, int64, error)
	// GetAll 获取所有关注
	GetAll(offset, limit int) ([]follow.Follow, int64, error)
	// Count 获取关注总数
	Count() (int64, error)
	// CountFollowers 获取粉丝总数
	CountFollowers(userID string) (int64, error)
	// CountFollowing 获取关注总数
	CountFollowing(userID string) (int64, error)
	// HasFollowed 检查是否已关注
	HasFollowed(userID, followID string) (bool, error)
}

type FollowRepositoryImpl struct{}

// NewFollowRepository 创建关注仓储
func NewFollowRepository() FollowRepository {
	return &FollowRepositoryImpl{}
}

// Create 创建关注关系
func (r *FollowRepositoryImpl) Create(f *follow.Follow) error {
	if err := database.DB.Create(f).Error; err != nil {
		return err
	}
	return nil
}

// Delete 删除关注关系
func (r *FollowRepositoryImpl) Delete(userID, followID string) error {
	if err := database.DB.Where("user_id = ? AND follow_id = ?", userID, followID).
		Delete(&follow.Follow{}).Error; err != nil {
		return err
	}
	return nil
}

// GetByID 根据ID获取关注
func (r *FollowRepositoryImpl) GetByID(id string) (*follow.Follow, error) {
	var f follow.Follow
	if err := database.DB.Where("id = ?", id).First(&f).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

// GetFollowers 获取粉丝列表
func (r *FollowRepositoryImpl) GetFollowers(userID string, offset, limit int) ([]follow.Follow, int64, error) {
	var followers []follow.Follow
	var count int64

	if err := database.DB.Where("follow_id = ?", userID).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := database.DB.
		Where("follow_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&followers).Error; err != nil {
		return nil, 0, err
	}

	return followers, count, nil
}

// GetFollowing 获取关注列表
func (r *FollowRepositoryImpl) GetFollowing(userID string, offset, limit int) ([]follow.Follow, int64, error) {
	var following []follow.Follow
	var count int64

	if err := database.DB.Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := database.DB.
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&following).Error; err != nil {
		return nil, 0, err
	}

	return following, count, nil
}

// GetAll 获取所有关注
func (r *FollowRepositoryImpl) GetAll(offset, limit int) ([]follow.Follow, int64, error) {
	var follows []follow.Follow
	var count int64

	if err := database.DB.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := database.DB.
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&follows).Error; err != nil {
		return nil, 0, err
	}

	return follows, count, nil
}

// Count 获取关注总数
func (r *FollowRepositoryImpl) Count() (int64, error) {
	var count int64
	if err := database.DB.Model(&follow.Follow{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountFollowers 获取粉丝总数
func (r *FollowRepositoryImpl) CountFollowers(userID string) (int64, error) {
	var count int64
	if err := database.DB.Model(&follow.Follow{}).Where("follow_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountFollowing 获取关注总数
func (r *FollowRepositoryImpl) CountFollowing(userID string) (int64, error) {
	var count int64
	if err := database.DB.Model(&follow.Follow{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// HasFollowed 检查是否已关注
func (r *FollowRepositoryImpl) HasFollowed(userID, followID string) (bool, error) {
	var count int64
	if err := database.DB.Model(&follow.Follow{}).
		Where("user_id = ? AND follow_id = ?", userID, followID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
