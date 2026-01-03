// Package services Topic业务逻辑服务
package services

import (
	"context"
	"fmt"
	"time"

	"GoHub-Service/app/cache"
	"GoHub-Service/app/models/topic"
	"GoHub-Service/app/repositories"
	apperrors "GoHub-Service/pkg/errors"
	"GoHub-Service/pkg/mapper"
	"GoHub-Service/pkg/paginator"
	"GoHub-Service/pkg/singleflight"

	"github.com/gin-gonic/gin"
)

// TopicService Topic服务
type TopicService struct {
	repo    repositories.TopicRepository
	cache   *cache.TopicCache
	sfGroup singleflight.Group                           // singleflight 防止缓存击穿
	mapper  mapper.Mapper[topic.Topic, TopicResponseDTO] // 使用泛型Mapper消除DTO转换重复
}

// NewTopicService 创建Topic服务实例
func NewTopicService() *TopicService {
	// 定义DTO转换函数（只需一次）
	converter := func(t *topic.Topic) *TopicResponseDTO {
		return &TopicResponseDTO{
			ID:            t.GetStringID(),
			Title:         t.Title,
			Body:          t.Body,
			CategoryID:    t.CategoryID,
			UserID:        t.UserID,
			LikeCount:     t.LikeCount,
			FavoriteCount: t.FavoriteCount,
			ViewCount:     t.ViewCount,
			CreatedAt:     t.CreatedAt,
			UpdatedAt:     t.UpdatedAt,
		}
	}

	return &TopicService{
		repo:   repositories.NewTopicRepository(),
		cache:  cache.NewTopicCache(),
		mapper: mapper.NewSimpleMapper(converter),
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
	LikeCount     int64     `json:"like_count"`
	FavoriteCount int64     `json:"favorite_count"`
	ViewCount     int64     `json:"view_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TopicListResponseDTO 话题列表响应DTO
type TopicListResponseDTO struct {
	Topics []TopicResponseDTO `json:"topics"`
	Paging *paginator.Paging  `json:"paging"`
}

// toResponseDTO 使用Mapper将Topic模型转换为响应DTO
// 优化：使用泛型Mapper消除重复代码
func (s *TopicService) toResponseDTO(t *topic.Topic) *TopicResponseDTO {
	return s.mapper.ToDTO(t)
}

// toResponseDTOList 使用Mapper将Topic模型列表转换为响应DTO列表
// 优化：使用泛型Mapper消除重复代码，自动优化内存拷贝
func (s *TopicService) toResponseDTOList(topics []topic.Topic) []TopicResponseDTO {
	return s.mapper.ToDTOList(topics)
}
		}
	}
	return dtos
}

// GetByID 根据ID获取话题（使用 singleflight 防止缓存击穿）
func (s *TopicService) GetByID(id string) (*TopicResponseDTO, *apperrors.AppError) {
	key := fmt.Sprintf("topic:%s", id)
	
	result, err := s.sfGroup.Do(key, func() (interface{}, error) {
		// 尝试从缓存获取
		if s.cache != nil {
			topicModel, err := s.cache.GetByID(context.Background(), id)
			if err == nil && topicModel != nil {
				return topicModel, nil
			}
		}

		// 从仓储获取
		topicModel, err := s.repo.GetByID(context.Background(), id)
		if err != nil {
			return nil, err
		}
		if topicModel == nil {
			return nil, apperrors.NotFoundError("话题").WithDetails(map[string]interface{}{"topic_id": id})
		}
		
		// 更新缓存
		if s.cache != nil {
			s.cache.Set(context.Background(), topicModel)
		}
		
		return topicModel, nil
	})
	
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return nil, appErr
		}
		return nil, apperrors.WrapError(err, "获取话题失败")
	}
	
	topicModel := result.(*topic.Topic)
	return s.toResponseDTO(topicModel), nil
}

// List 获取话题列表
func (s *TopicService) List(c *gin.Context, perPage int) (*TopicListResponseDTO, *apperrors.AppError) {
	if s.cache != nil {
		topics, found := s.cache.GetList(context.Background(), c)
		if found && len(topics) > 0 {
			data, pager, _ := s.repo.List(context.Background(), c, perPage)
			return &TopicListResponseDTO{
				Topics: s.toResponseDTOList(data),
				Paging: pager,
			}, nil
		}
	}

	data, pager, err := s.repo.List(context.Background(), c, perPage)
	if err != nil {
		return nil, apperrors.WrapError(err, "获取话题列表失败")
	}
	if len(data) > 0 && s.cache != nil {
		s.cache.SetList(context.Background(), c, data)
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
	if err := s.repo.Create(context.Background(), topicModel); err != nil {
		return nil, apperrors.WrapError(err, "创建话题失败")
	}
	if s.cache != nil {
		s.cache.ClearList(context.Background())
	}
	return s.toResponseDTO(topicModel), nil
}

// Update 更新话题
func (s *TopicService) Update(id string, dto TopicUpdateDTO) (*TopicResponseDTO, *apperrors.AppError) {
	topicModel, err := s.repo.GetByID(context.Background(), id)
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
	if err := s.repo.Update(context.Background(), topicModel); err != nil {
		return nil, apperrors.WrapError(err, "更新话题失败")
	}
	if s.cache != nil {
		s.cache.Delete(context.Background(), id)
		s.cache.ClearList(context.Background())
	}
	return s.toResponseDTO(topicModel), nil
}

// Delete 删除话题
func (s *TopicService) Delete(id string) *apperrors.AppError {
	err := s.repo.Delete(context.Background(), id)
	if err != nil {
		return apperrors.WrapError(err, "删除话题失败")
	}
	if s.cache != nil {
		s.cache.Delete(context.Background(), id)
		s.cache.ClearList(context.Background())
	}
	return nil
}

// CheckOwnership 检查用户是否拥有该话题
func (s *TopicService) CheckOwnership(topicID, userID string) (bool, *apperrors.AppError) {
	topicModel, err := s.repo.GetByID(context.Background(), topicID)
	if err != nil {
		return false, apperrors.WrapError(err, "检查话题所有权失败")
	}
	return topicModel.UserID == userID, nil
}
