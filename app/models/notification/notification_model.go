// Package notification 通知模型
package notification

import (
	"database/sql"
	"fmt"
	"time"

	"GoHub-Service/app/models"
	"GoHub-Service/pkg/database"

	"gorm.io/datatypes"
)

// Notification 通知记录
type Notification struct {
	models.BaseModel

	UserID  string         `json:"user_id,omitempty"`
	ActorID string         `json:"actor_id,omitempty"`
	Type    string         `json:"type,omitempty"`
	Data    datatypes.JSON `json:"data,omitempty"`
	ReadAt  sql.NullTime   `json:"read_at,omitempty"`

	models.CommonTimestampsField
}

// Create 创建通知
func (n *Notification) Create() {
	database.DB.Create(&n)
}

// MarkRead 标记已读
func (n *Notification) MarkRead() {
	now := time.Now()
	n.ReadAt = sql.NullTime{Time: now, Valid: true}
	database.DB.Model(&n).Update("read_at", n.ReadAt)
}

// GetID 实现Model接口
func (n *Notification) GetID() uint64 {
	return n.ID
}

// GetStringID 获取字符串ID
func (n *Notification) GetStringID() string {
	return fmt.Sprintf("%d", n.ID)
}
