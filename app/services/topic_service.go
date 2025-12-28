// Package services Topic业务逻辑服务
package services

import (
	"GoHub-Service/app/cache"
	"GoHub-Service/app/models/topic"
	"GoHub-Service/app/repositories"
	apperrors "GoHub-Service/pkg/errors"
	"GoHub-Service/pkg/paginator"

	"github.com/gin-gonic/gin"
)

// TopicService Topic服务
type TopicService struct{
	repo  repositories.TopicRepository
	cache *cache.TopicCache
}

// NewTopicService 创建Topic服务实例
func NewTopicService() *TopicService {
	return &TopicService{
		repo:  repositories.NewTopicRepository(),
		cache: cache.NewTopicCache(),
	}
}

// TopicCreateDTO 创建话题DTO
type TopicCreateDTO struct {
	Title      string
	Body       string
	CategoryID string
	UserID     string
}

// TopicUpdateDTO 更新话题DTO
type TopicUpdateDTO struct {
	Title      string
	Body       string
	CategoryID string
}

// GetByID 根据ID获取话题
func (s *TopicService) GetByID(id string) (*topic.Topic, error) {
	// 先从缓存获取
	topicModel, err := s.cache.GetByID(id)
	if err == nil && topicModel != nil {
		return topicModel, nil
	}
	
	// 缓存未命中，从数据库获取
	topicModel, err = s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if topicModel == nil {
		return nil, apperrors.NotFoundError("话题").WithDetails(map[string]interface{}{
			"topic_id": id,
		})
	}
	
	// 设置缓存
	s.cache.Set(topicModel)
	
	return topicModel, nil
}

// List 获取话题列表
func (s *TopicService) List(c *gin.Context, perPage int) ([]topic.Topic, *paginator.Paging, error) {
	// 先从缓存获取
	topics, found := s.cache.GetList(c)
	if found && len(topics) > 0 {
		// 简化处理，实际项目中应该缓存完整的分页信息
		data, pager, _ := s.repo.List(c, perPage)
		return data, pager, nil
	}
	
	// 缓存未命中，从数据库获取
	data, pager, err := s.repo.List(c, perPage)
	if err != nil {
		return nil, nil, err
	}
	
	// 设置缓存
	if len(data) > 0 {
		s.cache.SetList(c, data)
	}
	
	return data, pager, nil
}

// Create 创建话题
func (s *TopicService) Create(dto TopicCreateDTO) (*topic.Topic, error) {
	topicModel := &topic.Topic{
		Title:      dto.Title,
		Body:       dto.Body,
		CategoryID: dto.CategoryID,
		UserID:     dto.UserID,
	}

	if err := s.repo.Create(topicModel); err != nil {
		return nil, err
	}

	// 清除列表缓存
	s.cache.ClearList()

	return topicModel, nil
}

// Update 更新话题
func (s *TopicService) Update(id string, dto TopicUpdateDTO) (*topic.Topic, error) {
	topicModel, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if topicModel == nil {
		return nil, apperrors.NotFoundError("话题").WithDetails(map[string]interface{}{
			"topic_id": id,
		})
	}

	topicModel.Title = dto.Title
	topicModel.Body = dto.Body
	topicModel.CategoryID = dto.CategoryID

	if err := s.repo.Update(topicModel); err != nil {
		return nil, err
	}

	// 删除缓存
	s.cache.Delete(id)
	s.cache.ClearList()

	return topicModel, nil
}

// Delete 删除话题
func (s *TopicService) Delete(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	
	// 删除缓存
	s.cache.Delete(id)
	s.cache.ClearList()
	
	return nil
}

// CheckOwnership 检查用户是否拥有该话题
func (s *TopicService) CheckOwnership(topicID, userID string) (bool, error) {
	topicModel, err := s.repo.GetByID(topicID)
	if err != nil {
		return false, err
	}
	return topicModel.UserID == userID, nil
}

