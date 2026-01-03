// Package services 通知业务逻辑
package services

import (
	"context"
	"time"

	"GoHub-Service/app/models/notification"
	"GoHub-Service/app/repositories"
	apperrors "GoHub-Service/pkg/errors"
	"GoHub-Service/pkg/paginator"
	"GoHub-Service/pkg/resource"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NotificationService 通知服务
// 优化：使用GoRoutinePool管理批量通知的并发控制
type NotificationService struct {
	repo   repositories.NotificationRepository
	pool   *resource.GoRoutinePool // 使用goroutine池防止无限制创建goroutine
	logger *zap.Logger
}

// NewNotificationService 创建实例
// 优化：创建goroutine池，限制并发数为20
func NewNotificationService() *NotificationService {
	logger := zap.L()
	return &NotificationService{
		repo:   repositories.NewNotificationRepository(),
		pool:   resource.NewGoRoutinePool(20, logger),
		logger: logger,
	}
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

// BatchNotify 批量通知多个用户（使用goroutine池）
// 优化：使用GoRoutinePool限制并发，防止goroutine泄漏
func (s *NotificationService) BatchNotify(userIDs []string, actorID, typ string, data map[string]interface{}) *apperrors.AppError {
	if len(userIDs) == 0 {
		return nil
	}

	// 使用context控制超时
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	guard := resource.NewContextGuard(ctx, cancel, s.logger)
	defer guard.Release()

	// 使用goroutine池并发发送通知
	for _, userID := range userIDs {
		userID := userID // 捕获循环变量
		if err := s.pool.Submit(func() {
			// 检查context是否已取消
			select {
			case <-ctx.Done():
				s.logger.Warn("批量通知被取消", zap.String("user_id", userID))
				return
			default:
			}

			// 创建通知
			if err := s.repo.Create(userID, actorID, typ, data); err != nil {
				s.logger.Error("创建通知失败",
					zap.Error(err),
					zap.String("user_id", userID),
					zap.String("type", typ),
				)
			}
		}); err != nil {
			s.logger.Error("提交通知任务失败", zap.Error(err))
			return apperrors.InternalError("批量通知失败", err)
		}
	}

	guard.Cancel() // 操作完成，提前取消
	return nil
}

// Shutdown 优雅关闭服务
func (s *NotificationService) Shutdown(timeout time.Duration) error {
	s.logger.Info("正在关闭NotificationService...")
	if err := s.pool.Shutdown(timeout); err != nil {
		s.logger.Error("关闭goroutine池失败", zap.Error(err))
		return err
	}
	s.logger.Info("NotificationService已关闭")
	return nil
}
