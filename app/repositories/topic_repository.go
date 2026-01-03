// Package repositories 话题数据访问层
package repositories

import (
	"context"
	"GoHub-Service/app/models/topic"
	"GoHub-Service/pkg/paginator"
	"GoHub-Service/pkg/database"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// TopicRepository 话题仓储接口
type TopicRepository interface {
	GetByID(ctx context.Context, id string) (*topic.Topic, error)
	List(ctx context.Context, c *gin.Context, perPage int) ([]topic.Topic, *paginator.Paging, error)
	Create(ctx context.Context, topic *topic.Topic) error
	Update(ctx context.Context, topic *topic.Topic) error
	Delete(ctx context.Context, id string) error
	BatchCreate(ctx context.Context, topics []topic.Topic) error
	BatchDelete(ctx context.Context, ids []string) error
}

// topicRepository 话题仓储实现
type topicRepository struct{}

// BatchCreate 批量创建话题（事务包裹）
func (r *topicRepository) BatchCreate(ctx context.Context, topics []topic.Topic) error {
	return database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&topics).Error; err != nil {
			return err
		}
		return nil
	})
}

// BatchDelete 批量删除话题（事务包裹）
func (r *topicRepository) BatchDelete(ctx context.Context, ids []string) error {
	return database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id IN ?", ids).Delete(&topic.Topic{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// NewTopicRepository 创建话题仓储实例
func NewTopicRepository() TopicRepository {
	return &topicRepository{}
}

// GetByID 根据ID获取话题
func (r *topicRepository) GetByID(ctx context.Context, id string) (*topic.Topic, error) {
	topicModel := topic.Get(id)
	if topicModel.ID == 0 {
		return nil, nil
	}
	return &topicModel, nil
}

// List 获取话题列表
func (r *topicRepository) List(ctx context.Context, c *gin.Context, perPage int) ([]topic.Topic, *paginator.Paging, error) {
	data, pager := topic.Paginate(c, perPage)
	return data, &pager, nil
}

// Create 创建话题
func (r *topicRepository) Create(ctx context.Context, t *topic.Topic) error {
	t.Create()
	if t.ID == 0 {
		return ErrCreateFailed
	}
	return nil
}

// Update 更新话题
func (r *topicRepository) Update(ctx context.Context, t *topic.Topic) error {
	rowsAffected := t.Save()
	if rowsAffected == 0 {
		return ErrUpdateFailed
	}
	return nil
}

// Delete 删除话题
func (r *topicRepository) Delete(ctx context.Context, id string) error {
	topicModel := topic.Get(id)
	if topicModel.ID == 0 {
		return ErrNotFound
	}

	rowsAffected := topicModel.Delete()
	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return nil
}
