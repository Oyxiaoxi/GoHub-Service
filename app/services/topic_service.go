// Package services Topic业务逻辑服务
package services

import (
	"GoHub-Service/app/cache"
	"GoHub-Service/app/models/topic"
	"GoHub-Service/app/repositories"
	apperrors "GoHub-Service/pkg/errors"
	"GoHub-Service/pkg/paginator"
	"time"

	"github.com/gin-gonic/gin"
)

// TopicService Topic服务
type TopicService struct {
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
	Title      string `json:"title" binding:"required,min=3,max=255"`
	Body       string `json:"body" binding:"required"`
	CategoryID string `json:"category_id" binding:"required"`
	UserID     string `json:"user_id" binding:"required"`
}

// TopicUpdateDTO 更新话题DTO
type TopicUpdateDTO struct {
	Title      *string `json:"title,omitempty" binding:"omitempty,min=3,max=255"`
	Body       *string `json:"body,omitempty"`
	CategoryID *string `json:"category_id,omitempty"`
}

// TopicResponseDTO 话题响应DTO
type TopicResponseDTO struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	CategoryID string    `json:"category_id"`
	UserID     string    `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TopicListResponseDTO 话题列表响应DTO
type TopicListResponseDTO struct {
	Topics []TopicResponseDTO `json:"topics"`
	Paging *paginator.Paging  `json:"paging"`
}

// toResponseDTO 将Topic模型转换为响应DTO
func (s *TopicService) toResponseDTO(t *topic.Topic) *TopicResponseDTO {
	return &TopicResponseDTO{
		ID:         t.GetStringID(),
		Title:      t.Title,
		Body:       t.Body,
		CategoryID: t.CategoryID,
		UserID:     t.UserID,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}

// toResponseDTOList 将Topic模型列表转换为响应DTO列表
func (s *TopicService) toResponseDTOList(topics []topic.Topic) []TopicResponseDTO {
	dtos := make([]TopicResponseDTO, len(topics))
	for i, t := range topics {
		dtos[i] = TopicResponseDTO{
			ID:         t.GetStringID(),
			Title:      t.Title,
			Body:       t.Body,
			CategoryID: t.CategoryID,
			UserID:     t.UserID,
			CreatedAt:  t.CreatedAt,
			UpdatedAt:  t.UpdatedAt,
		}
	}
	return dtos
}

// GetByID 根据ID获取话题
func (s *TopicService) GetByID(id string) (*TopicResponseDTO, *apperrors.AppError) {
	if s.cache != nil {
		topicModel, err := s.cache.GetByID(id)
		if err == nil && topicModel != nil {
			return s.toResponseDTO(topicModel), nil
		}
	}

	topicModel, err := s.repo.GetByID(id)
	if err != nil {
		return nil, apperrors.WrapError(err, "获取话题失败")
	}
	if topicModel == nil {
		return nil, apperrors.NotFoundError("话题").WithDetails(map[string]interface{}{"topic_id": id})
	}
	if s.cache != nil {
		s.cache.Set(topicModel)
	}
	return s.toResponseDTO(topicModel), nil
}

// List 获取话题列表
func (s *TopicService) List(c *gin.Context, perPage int) (*TopicListResponseDTO, *apperrors.AppError) {
	if s.cache != nil {
		topics, found := s.cache.GetList(c)
		if found && len(topics) > 0 {
			data, pager, _ := s.repo.List(c, perPage)
			return &TopicListResponseDTO{
				Topics: s.toResponseDTOList(data),
				Paging: pager,
			}, nil
		}
	}

	data, pager, err := s.repo.List(c, perPage)
	if err != nil {
		return nil, apperrors.WrapError(err, "获取话题列表失败")
	}
	if len(data) > 0 && s.cache != nil {
		s.cache.SetList(c, data)
	}
	return &TopicListResponseDTO{
		Topics: s.toResponseDTOList(data),
		Paging: pager,
	}, nil
}

// Create 创建话题
func (s *TopicService) Create(dto TopicCreateDTO) (*TopicResponseDTO, *apperrors.AppError) {
	topicModel := &topic.Topic{
		Title:      dto.Title,
		Body:       dto.Body,
		CategoryID: dto.CategoryID,
		UserID:     dto.UserID,
	}
	if err := s.repo.Create(topicModel); err != nil {
		return nil, apperrors.WrapError(err, "创建话题失败")
	}
	if s.cache != nil {
		s.cache.ClearList()
	}
	return s.toResponseDTO(topicModel), nil
}

// Update 更新话题
func (s *TopicService) Update(id string, dto TopicUpdateDTO) (*TopicResponseDTO, *apperrors.AppError) {
	topicModel, err := s.repo.GetByID(id)
	if err != nil {
		return nil, apperrors.WrapError(err, "获取话题失败")
	}
	if topicModel == nil {
		return nil, apperrors.NotFoundError("话题").WithDetails(map[string]interface{}{"topic_id": id})
	}
	if dto.Title != nil {
		topicModel.Title = *dto.Title
	}
	if dto.Body != nil {
		topicModel.Body = *dto.Body
	}
	if dto.CategoryID != nil {
		topicModel.CategoryID = *dto.CategoryID
	}
	if err := s.repo.Update(topicModel); err != nil {
		return nil, apperrors.WrapError(err, "更新话题失败")
	}
	if s.cache != nil {
		s.cache.Delete(id)
		s.cache.ClearList()
	}
	return s.toResponseDTO(topicModel), nil
}

// Delete 删除话题
func (s *TopicService) Delete(id string) *apperrors.AppError {
	err := s.repo.Delete(id)
	if err != nil {
		return apperrors.WrapError(err, "删除话题失败")
	}
	if s.cache != nil {
		s.cache.Delete(id)
		s.cache.ClearList()
	}
	return nil
}

// CheckOwnership 检查用户是否拥有该话题
func (s *TopicService) CheckOwnership(topicID, userID string) (bool, *apperrors.AppError) {
	topicModel, err := s.repo.GetByID(topicID)
	if err != nil {
		return false, apperrors.WrapError(err, "检查话题所有权失败")
	}
	return topicModel.UserID == userID, nil
}
