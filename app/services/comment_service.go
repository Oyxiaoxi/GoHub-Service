// Package services Comment业务逻辑服务
package services

import (
	"context"
	"fmt"
	"time"

	"GoHub-Service/app/cache"
	"GoHub-Service/app/models/comment"
	"GoHub-Service/app/repositories"
	apperrors "GoHub-Service/pkg/errors"
	"GoHub-Service/pkg/mapper"
	"GoHub-Service/pkg/paginator"
	"GoHub-Service/pkg/singleflight"

	"github.com/gin-gonic/gin"
)

// CommentService 评论服务
type CommentService struct {
	repo      repositories.CommentRepository
	cache     *cache.CommentCache
	notifSvc  *NotificationService
	topicRepo repositories.TopicRepository
	sfGroup   singleflight.Group                               // singleflight 防止缓存击穿
	mapper    mapper.Mapper[comment.Comment, CommentResponseDTO] // 使用泛型Mapper消除DTO转换重复
}

// NewCommentService 创建评论服务实例
func NewCommentService() *CommentService {
	// 定义DTO转换函数（只需一次）
	converter := func(c *comment.Comment) *CommentResponseDTO {
		return &CommentResponseDTO{
			ID:        c.GetStringID(),
			TopicID:   c.TopicID,
			UserID:    c.UserID,
			Content:   c.Content,
			ParentID:  c.ParentID,
			LikeCount: c.LikeCount,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}
	}

	return &CommentService{
		repo:      repositories.NewCommentRepository(),
		cache:     cache.NewCommentCache(),
		notifSvc:  NewNotificationService(),
		topicRepo: repositories.NewTopicRepository(),
		mapper:    mapper.NewSimpleMapper(converter),
	}
}

// CommentCreateDTO 创建评论DTO
type CommentCreateDTO struct {
	TopicID  string `json:"topic_id" binding:"required"`
	UserID   string `json:"user_id" binding:"required"`
	Content  string `json:"content" binding:"required,min=1,max=1000"`
	ParentID string `json:"parent_id,omitempty"`
}

// CommentUpdateDTO 更新评论DTO
type CommentUpdateDTO struct {
	Content *string `json:"content,omitempty" binding:"omitempty,min=1,max=1000"`
}

// CommentResponseDTO 评论响应DTO
type CommentResponseDTO struct {
	ID        string    `json:"id"`
	TopicID   string    `json:"topic_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	ParentID  string    `json:"parent_id"`
	LikeCount int64     `json:"like_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CommentListResponseDTO 评论列表响应DTO
type CommentListResponseDTO struct {
	Comments []CommentResponseDTO `json:"comments"`
	Paging   *paginator.Paging    `json:"paging"`
}

// toResponseDTO 使用Mapper将Comment模型转换为响应DTO
// 优化：使用泛型Mapper消除重复代码
func (s *CommentService) toResponseDTO(c *comment.Comment) *CommentResponseDTO {
	return s.mapper.ToDTO(c)
}

// toResponseDTOList 使用Mapper将Comment模型列表转换为响应DTO列表
// 优化：使用泛型Mapper消除重复代码，自动优化内存拷贝
func (s *CommentService) toResponseDTOList(comments []comment.Comment) []CommentResponseDTO {
	return s.mapper.ToDTOList(comments)
}

// GetByID 根据ID获取评论（使用 singleflight 防止缓存击穿）
func (s *CommentService) GetByID(ctx context.Context, id string) (*CommentResponseDTO, *apperrors.AppError) {
	// 使用 singleflight 确保同一时间只有一个请求去数据库查询
	key := fmt.Sprintf("comment:%s", id)
	
	result, err := s.sfGroup.Do(key, func() (interface{}, error) {
		// 尝试从缓存获取
		if s.cache != nil {
			commentModel, err := s.cache.GetByID(ctx, id)
			if err == nil && commentModel != nil {
				return commentModel, nil
			}
		}

		// 从仓储获取
		commentModel, err := s.repo.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		if commentModel == nil {
			return nil, apperrors.NotFoundError("评论")
		}

		// 更新缓存
		if s.cache != nil {
			s.cache.Set(ctx, commentModel)
		}

		return commentModel, nil
	})

	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return nil, appErr
		}
		return nil, apperrors.DatabaseError("获取评论", err)
	}

	commentModel := result.(*comment.Comment)
	return s.toResponseDTO(commentModel), nil
}

// List 获取评论列表
func (s *CommentService) List(ctx context.Context, c *gin.Context, perPage int) (*CommentListResponseDTO, *apperrors.AppError) {
	comments, paging, err := s.repo.List(ctx, c, perPage)
	if err != nil {
		return nil, apperrors.DatabaseError("获取评论列表", err)
	}

	return &CommentListResponseDTO{
		Comments: s.toResponseDTOList(comments),
		Paging:   paging,
	}, nil
}

// ListByTopicID 获取指定话题的评论列表
func (s *CommentService) ListByTopicID(ctx context.Context, c *gin.Context, topicID string, perPage int) (*CommentListResponseDTO, *apperrors.AppError) {
	// 尝试从缓存获取
	if s.cache != nil {
		comments, err := s.cache.GetByTopicID(ctx, topicID)
		if err == nil && comments != nil {
			// 缓存命中，但仍需分页处理
			// 这里简化处理，实际应用中可以优化缓存分页
		}
	}

	comments, paging, err := s.repo.ListByTopicID(ctx, c, topicID, perPage)
	if err != nil {
		return nil, apperrors.DatabaseError("获取话题评论列表", err)
	}

	return &CommentListResponseDTO{
		Comments: s.toResponseDTOList(comments),
		Paging:   paging,
	}, nil
}

// ListByUserID 获取指定用户的评论列表
func (s *CommentService) ListByUserID(ctx context.Context, c *gin.Context, userID string, perPage int) (*CommentListResponseDTO, *apperrors.AppError) {
	comments, paging, err := s.repo.ListByUserID(ctx, c, userID, perPage)
	if err != nil {
		return nil, apperrors.DatabaseError("获取用户评论列表", err)
	}

	return &CommentListResponseDTO{
		Comments: s.toResponseDTOList(comments),
		Paging:   paging,
	}, nil
}

// ListReplies 获取评论的回复列表
func (s *CommentService) ListReplies(ctx context.Context, c *gin.Context, parentID string, perPage int) (*CommentListResponseDTO, *apperrors.AppError) {
	comments, paging, err := s.repo.ListReplies(ctx, c, parentID, perPage)
	if err != nil {
		return nil, apperrors.DatabaseError("获取评论回复列表", err)
	}

	return &CommentListResponseDTO{
		Comments: s.toResponseDTOList(comments),
		Paging:   paging,
	}, nil
}

// Create 创建评论
func (s *CommentService) Create(ctx context.Context, dto *CommentCreateDTO) (*CommentResponseDTO, *apperrors.AppError) {
	// 创建评论模型
	commentModel := &comment.Comment{
		TopicID:  dto.TopicID,
		UserID:   dto.UserID,
		Content:  dto.Content,
		ParentID: dto.ParentID,
	}

	// 如果没有指定父评论，设置为0（顶级评论）
	if commentModel.ParentID == "" {
		commentModel.ParentID = "0"
	}

	// 保存到数据库
	if err := s.repo.Create(ctx, commentModel); err != nil {
		return nil, apperrors.DatabaseError("创建评论", err)
	}

	// 发送通知：话题作者、父评论作者
	if s.notifSvc != nil {
		// 通知话题作者
		if topicModel, err := s.topicRepo.GetByID(ctx, dto.TopicID); err == nil && topicModel != nil {
			if topicModel.UserID != "" && topicModel.UserID != dto.UserID {
				_ = s.notifSvc.Notify(topicModel.UserID, dto.UserID, "comment_created", map[string]interface{}{"topic_id": dto.TopicID, "comment_id": commentModel.GetStringID()})
			}
		}
		// 通知父评论作者
		if dto.ParentID != "" && dto.ParentID != "0" {
			if parent, err := s.repo.GetByID(ctx, dto.ParentID); err == nil && parent != nil {
				if parent.UserID != "" && parent.UserID != dto.UserID {
					_ = s.notifSvc.Notify(parent.UserID, dto.UserID, "comment_replied", map[string]interface{}{"topic_id": dto.TopicID, "comment_id": commentModel.GetStringID(), "parent_id": dto.ParentID})
				}
			}
		}
	}

	// 清除相关缓存
	if s.cache != nil {
		s.cache.InvalidateByTopicID(ctx, dto.TopicID)
	}

	return s.toResponseDTO(commentModel), nil
}

// Update 更新评论
func (s *CommentService) Update(ctx context.Context, id string, dto *CommentUpdateDTO) (*CommentResponseDTO, *apperrors.AppError) {
	// 获取评论
	commentModel, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.DatabaseError("获取评论", err)
	}
	if commentModel == nil {
		return nil, apperrors.NotFoundError("评论")
	}

	// 更新字段
	if dto.Content != nil {
		commentModel.Content = *dto.Content
	}

	// 保存更新
	if err := s.repo.Update(ctx, commentModel); err != nil {
		return nil, apperrors.DatabaseError("更新评论", err)
	}

	// 清除缓存
	if s.cache != nil {
		s.cache.Invalidate(ctx, id)
		s.cache.InvalidateByTopicID(ctx, commentModel.TopicID)
	}

	return s.toResponseDTO(commentModel), nil
}

// Delete 删除评论
func (s *CommentService) Delete(ctx context.Context, id string) *apperrors.AppError {
	// 获取评论信息用于清除缓存
	commentModel, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperrors.DatabaseError("获取评论", err)
	}
	if commentModel == nil {
		return apperrors.NotFoundError("评论")
	}

	// 删除评论
	if err := s.repo.Delete(ctx, id); err != nil {
		return apperrors.DatabaseError("删除评论", err)
	}

	// 清除缓存
	if s.cache != nil {
		s.cache.Invalidate(ctx, id)
		s.cache.InvalidateByTopicID(ctx, commentModel.TopicID)
	}

	return nil
}

// LikeComment 点赞评论
func (s *CommentService) LikeComment(ctx context.Context, id string) *apperrors.AppError {
	if err := s.repo.IncrementLikeCount(ctx, id); err != nil {
		return apperrors.DatabaseError("点赞评论", err)
	}

	// 清除缓存
	if s.cache != nil {
		s.cache.Invalidate(ctx, id)
	}

	return nil
}

// UnlikeComment 取消点赞
func (s *CommentService) UnlikeComment(ctx context.Context, id string) *apperrors.AppError {
	if err := s.repo.DecrementLikeCount(ctx, id); err != nil {
		return apperrors.DatabaseError("取消点赞评论", err)
	}

	// 清除缓存
	if s.cache != nil {
		s.cache.Invalidate(ctx, id)
	}

	return nil
}

// CountByTopicID 统计话题的评论数
func (s *CommentService) CountByTopicID(ctx context.Context, topicID string) (int64, *apperrors.AppError) {
	count, err := s.repo.CountByTopicID(ctx, topicID)
	if err != nil {
		return 0, apperrors.DatabaseError("统计评论数", err)
	}
	return count, nil
}
