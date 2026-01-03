// Package repositories 评论数据访问层
package repositories

import (
	"GoHub-Service/app/models/comment"
	"GoHub-Service/pkg/database"
	"GoHub-Service/pkg/paginator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CommentRepository 评论仓储接口
type CommentRepository interface {
	GetByID(id string) (*comment.Comment, error)
	List(c *gin.Context, perPage int) ([]comment.Comment, *paginator.Paging, error)
	ListByTopicID(c *gin.Context, topicID string, perPage int) ([]comment.Comment, *paginator.Paging, error)
	ListByUserID(c *gin.Context, userID string, perPage int) ([]comment.Comment, *paginator.Paging, error)
	ListReplies(c *gin.Context, parentID string, perPage int) ([]comment.Comment, *paginator.Paging, error)
	Create(comment *comment.Comment) error
	Update(comment *comment.Comment) error
	Delete(id string) error
	BatchCreate(comments []comment.Comment) error
	BatchDelete(ids []string) error
	IncrementLikeCount(id string) error
	DecrementLikeCount(id string) error
	CountByTopicID(topicID string) (int64, error)
}

// commentRepository 评论仓储实现
type commentRepository struct{}

// NewCommentRepository 创建评论仓储实例
func NewCommentRepository() CommentRepository {
	return &commentRepository{}
}

// GetByID 根据ID获取评论
func (r *commentRepository) GetByID(id string) (*comment.Comment, error) {
	var commentModel comment.Comment
	if err := database.DB.
		Select("id", "topic_id", "user_id", "content", "parent_id", "like_count", "created_at", "updated_at").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email", "avatar")
		}).
		Preload("Topic", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "title", "user_id", "category_id")
		}).
		First(&commentModel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &commentModel, nil
}

// List 获取评论列表
func (r *commentRepository) List(c *gin.Context, perPage int) ([]comment.Comment, *paginator.Paging, error) {
	var comments []comment.Comment
	query := database.DB.Model(&comment.Comment{}).
		Select("id", "topic_id", "user_id", "content", "parent_id", "like_count", "created_at", "updated_at").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email", "avatar")
		}).
		Preload("Topic", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "title", "user_id", "category_id")
		}).
		Order("created_at DESC")

	paging := paginator.Paginate(
		c,
		query,
		&comments,
		"/api/v1/comments",
		perPage,
	)

	return comments, &paging, nil
}

// ListByTopicID 获取指定话题的评论列表
func (r *commentRepository) ListByTopicID(c *gin.Context, topicID string, perPage int) ([]comment.Comment, *paginator.Paging, error) {
	var comments []comment.Comment
	query := database.DB.Model(&comment.Comment{}).
		Select("id", "topic_id", "user_id", "content", "parent_id", "like_count", "created_at", "updated_at").
		Where("topic_id = ? AND parent_id = ?", topicID, "0").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email", "avatar")
		}).
		Order("created_at DESC")

	paging := paginator.Paginate(
		c,
		query,
		&comments,
		"/api/v1/topics/"+topicID+"/comments",
		perPage,
	)

	return comments, &paging, nil
}

// ListByUserID 获取指定用户的评论列表
func (r *commentRepository) ListByUserID(c *gin.Context, userID string, perPage int) ([]comment.Comment, *paginator.Paging, error) {
	var comments []comment.Comment
	query := database.DB.Model(&comment.Comment{}).
		Select("id", "topic_id", "user_id", "content", "parent_id", "like_count", "created_at", "updated_at").
		Where("user_id = ?", userID).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email", "avatar")
		}).
		Preload("Topic", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "title", "user_id", "category_id")
		}).
		Order("created_at DESC")

	paging := paginator.Paginate(
		c,
		query,
		&comments,
		"/api/v1/users/"+userID+"/comments",
		perPage,
	)

	return comments, &paging, nil
}

// ListReplies 获取评论的回复列表
func (r *commentRepository) ListReplies(c *gin.Context, parentID string, perPage int) ([]comment.Comment, *paginator.Paging, error) {
	var comments []comment.Comment
	query := database.DB.Model(&comment.Comment{}).
		Select("id", "topic_id", "user_id", "content", "parent_id", "like_count", "created_at", "updated_at").
		Where("parent_id = ?", parentID).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email", "avatar")
		}).
		Order("created_at ASC")

	paging := paginator.Paginate(
		c,
		query,
		&comments,
		"/api/v1/comments/"+parentID+"/replies",
		perPage,
	)

	return comments, &paging, nil
}

// Create 创建评论
func (r *commentRepository) Create(c *comment.Comment) error {
	c.Create()
	if c.ID == 0 {
		return ErrCreateFailed
	}
	return nil
}

// Update 更新评论
func (r *commentRepository) Update(c *comment.Comment) error {
	rowsAffected := c.Save()
	if rowsAffected == 0 {
		return ErrUpdateFailed
	}
	return nil
}

// Delete 删除评论
func (r *commentRepository) Delete(id string) error {
	var commentModel comment.Comment
	if err := database.DB.First(&commentModel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return err
	}

	rowsAffected := commentModel.Delete()
	if rowsAffected == 0 {
		return ErrDeleteFailed
	}
	return nil
}

// IncrementLikeCount 增加点赞数
func (r *commentRepository) IncrementLikeCount(id string) error {
	return database.DB.Model(&comment.Comment{}).
		Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
}

// DecrementLikeCount 减少点赞数
func (r *commentRepository) DecrementLikeCount(id string) error {
	return database.DB.Model(&comment.Comment{}).
		Where("id = ? AND like_count > ?", id, 0).
		UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error
}

// CountByTopicID 统计话题的评论数
func (r *commentRepository) CountByTopicID(topicID string) (int64, error) {
	var count int64
	err := database.DB.Model(&comment.Comment{}).Where("topic_id = ?", topicID).Count(&count).Error
	return count, err
}

// BatchCreate 批量创建评论（使用事务和批量插入优化）
func (r *commentRepository) BatchCreate(comments []comment.Comment) error {
	if len(comments) == 0 {
		return nil
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 使用 CreateInBatches 批量插入，每批100条
		if err := tx.CreateInBatches(&comments, 100).Error; err != nil {
			return err
		}
		return nil
	})
}

// BatchDelete 批量删除评论（使用事务）
func (r *commentRepository) BatchDelete(ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id IN ?", ids).Delete(&comment.Comment{}).Error; err != nil {
			return err
		}
		return nil
	})
}
