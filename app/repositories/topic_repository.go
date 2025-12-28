// Package repositories 话题数据访问层
package repositories

import (
	"GoHub-Service/app/models/topic"
	"GoHub-Service/pkg/paginator"

	"github.com/gin-gonic/gin"
)

// TopicRepository 话题仓储接口
type TopicRepository interface {
	GetByID(id string) (*topic.Topic, error)
	List(c *gin.Context, perPage int) ([]topic.Topic, *paginator.Paging, error)
	Create(topic *topic.Topic) error
	Update(topic *topic.Topic) error
	Delete(id string) error
}

// topicRepository 话题仓储实现
type topicRepository struct{}

// NewTopicRepository 创建话题仓储实例
func NewTopicRepository() TopicRepository {
	return &topicRepository{}
}

// GetByID 根据ID获取话题
func (r *topicRepository) GetByID(id string) (*topic.Topic, error) {
	topicModel := topic.Get(id)
	if topicModel.ID == 0 {
		return nil, nil
	}
	return &topicModel, nil
}

// List 获取话题列表
func (r *topicRepository) List(c *gin.Context, perPage int) ([]topic.Topic, *paginator.Paging, error) {
	data, pager := topic.Paginate(c, perPage)
	return data, &pager, nil
}

// Create 创建话题
func (r *topicRepository) Create(t *topic.Topic) error {
	t.Create()
	if t.ID == 0 {
		return ErrCreateFailed
	}
	return nil
}

// Update 更新话题
func (r *topicRepository) Update(t *topic.Topic) error {
	rowsAffected := t.Save()
	if rowsAffected == 0 {
		return ErrUpdateFailed
	}
	return nil
}

// Delete 删除话题
func (r *topicRepository) Delete(id string) error {
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
