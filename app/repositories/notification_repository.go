// Package repositories 通知数据访问层
package repositories

import (
	"GoHub-Service/app/models/notification"
	"GoHub-Service/pkg/database"
	"GoHub-Service/pkg/paginator"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// NotificationRepository 通知仓储接口
type NotificationRepository interface {
	Create(userID, actorID, typ string, data map[string]interface{}) error
	ListByUser(c *gin.Context, userID string, perPage int) ([]notification.Notification, *paginator.Paging, error)
	MarkRead(id, userID string) error
	MarkAllRead(userID string) error
}

type notificationRepository struct{}

// NewNotificationRepository 创建通知仓储实例
func NewNotificationRepository() NotificationRepository {
	return &notificationRepository{}
}

// Create 创建通知记录
func (r *notificationRepository) Create(userID, actorID, typ string, data map[string]interface{}) error {
	jsonData, _ := json.Marshal(data)
	record := notification.Notification{
		UserID:  userID,
		ActorID: actorID,
		Type:    typ,
		Data:    datatypes.JSON(jsonData),
	}
	return database.DB.Create(&record).Error
}

// ListByUser 获取用户通知列表
func (r *notificationRepository) ListByUser(c *gin.Context, userID string, perPage int) ([]notification.Notification, *paginator.Paging, error) {
	var notifications []notification.Notification
	query := database.DB.Model(&notification.Notification{}).
		Where("user_id = ?", userID).
		Order("created_at DESC")

	paging := paginator.Paginate(c, query, &notifications, "/api/v1/notifications", perPage)
	return notifications, &paging, nil
}

// MarkRead 标记单条通知已读（限制用户）
func (r *notificationRepository) MarkRead(id, userID string) error {
	readAt := sql.NullTime{Time: time.Now(), Valid: true}
	return database.DB.Model(&notification.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("read_at", readAt).Error
}

// MarkAllRead 标记用户全部通知已读
func (r *notificationRepository) MarkAllRead(userID string) error {
	return database.DB.Model(&notification.Notification{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Update("read_at", gorm.Expr("NOW()")).Error
}
