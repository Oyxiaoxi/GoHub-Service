// Package services 通知业务逻辑
package services

import (
	"GoHub-Service/app/models/notification"
	"GoHub-Service/app/repositories"
	apperrors "GoHub-Service/pkg/errors"
	"GoHub-Service/pkg/paginator"

	"github.com/gin-gonic/gin"
)

// NotificationService 通知服务
type NotificationService struct {
	repo repositories.NotificationRepository
}

// NewNotificationService 创建实例
func NewNotificationService() *NotificationService {
	return &NotificationService{repo: repositories.NewNotificationRepository()}
}

// Notify 创建通知
func (s *NotificationService) Notify(userID, actorID, typ string, data map[string]interface{}) *apperrors.AppError {
	if userID == "" {
		return nil
	}
	if err := s.repo.Create(userID, actorID, typ, data); err != nil {
		return apperrors.DatabaseError("创建通知", err)
	}
	return nil
}

// List 获取通知列表
func (s *NotificationService) List(c *gin.Context, userID string, perPage int) ([]notification.Notification, *paginator.Paging, *apperrors.AppError) {
	list, paging, err := s.repo.ListByUser(c, userID, perPage)
	if err != nil {
		return nil, nil, apperrors.DatabaseError("获取通知列表", err)
	}
	return list, paging, nil
}

// MarkRead 标记通知已读
func (s *NotificationService) MarkRead(id, userID string) *apperrors.AppError {
	if err := s.repo.MarkRead(id, userID); err != nil {
		return apperrors.DatabaseError("标记通知已读", err)
	}
	return nil
}

// MarkAllRead 标记全部为已读
func (s *NotificationService) MarkAllRead(userID string) *apperrors.AppError {
	if err := s.repo.MarkAllRead(userID); err != nil {
		return apperrors.DatabaseError("标记全部通知已读", err)
	}
	return nil
}

// ToResponse 简单返回响应对象
func (s *NotificationService) ToResponse(n notification.Notification) map[string]interface{} {
	return map[string]interface{}{
		"id":         n.GetStringID(),
		"user_id":    n.UserID,
		"actor_id":   n.ActorID,
		"type":       n.Type,
		"data":       n.Data,
		"read_at":    n.ReadAt,
		"created_at": n.CreatedAt,
		"updated_at": n.UpdatedAt,
	}
}

// ToResponseList 转换列表
func (s *NotificationService) ToResponseList(list []notification.Notification) []map[string]interface{} {
	res := make([]map[string]interface{}, len(list))
	for i, n := range list {
		res[i] = s.ToResponse(n)
	}
	return res
}
