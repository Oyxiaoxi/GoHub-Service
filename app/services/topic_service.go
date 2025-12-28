// Package services Topic业务逻辑服务
package services

import (
	"GoHub-Service/app/models/topic"
	apperrors "GoHub-Service/pkg/errors"
	"GoHub-Service/pkg/paginator"

	"github.com/gin-gonic/gin"
)

// TopicService Topic服务
type TopicService struct{}

// NewTopicService 创建Topic服务实例
func NewTopicService() *TopicService {
	return &TopicService{}
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
	topicModel := topic.Get(id)
	if topicModel.ID == 0 {
		return nil, apperrors.NotFoundError("话题")
	}
	return &topicModel, nil
}

// List 获取话题列表（分页）
func (s *TopicService) List(c *gin.Context, perPage int) ([]topic.Topic, paginator.Paging, error) {
	data, paging := topic.Paginate(c, perPage)
	return data, paging, nil
}

// Create 创建话题
func (s *TopicService) Create(dto TopicCreateDTO) (*topic.Topic, error) {
	// 验证分类是否存在
	// 这里可以添加业务逻辑验证
	
	topicModel := topic.Topic{
		Title:      dto.Title,
		Body:       dto.Body,
		CategoryID: dto.CategoryID,
		UserID:     dto.UserID,
	}
	
	topicModel.Create()
	
	if topicModel.ID == 0 {
		return nil, apperrors.DatabaseError("创建话题", nil)
	}
	
	return &topicModel, nil
}

// Update 更新话题
func (s *TopicService) Update(id string, dto TopicUpdateDTO) (*topic.Topic, error) {
	// 获取话题
	topicModel, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	
	// 更新字段
	topicModel.Title = dto.Title
	topicModel.Body = dto.Body
	topicModel.CategoryID = dto.CategoryID
	
	// 保存
	rowsAffected := topicModel.Save()
	if rowsAffected == 0 {
		return nil, apperrors.DatabaseError("更新话题", nil)
	}
	
	return topicModel, nil
}

// Delete 删除话题
func (s *TopicService) Delete(id string) error {
	topicModel, err := s.GetByID(id)
	if err != nil {
		return err
	}
	
	rowsAffected := topicModel.Delete()
	if rowsAffected == 0 {
		return apperrors.DatabaseError("删除话题", nil)
	}
	
	return nil
}

// CheckOwnership 检查用户是否拥有话题
func (s *TopicService) CheckOwnership(topicID string, userID string) (bool, error) {
	topicModel, err := s.GetByID(topicID)
	if err != nil {
		return false, err
	}
	
	return topicModel.UserID == userID, nil
}
