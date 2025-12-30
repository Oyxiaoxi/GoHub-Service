package message

import (
	"time"

	"GoHub-Service/app/models"
	"GoHub-Service/pkg/database"
)

// Message 私信消息模型
type Message struct {
	models.BaseModel

	ConversationID string     `json:"conversation_id,omitempty"`
	SenderID       string     `json:"sender_id,omitempty"`
	ReceiverID     string     `json:"receiver_id,omitempty"`
	Body           string     `json:"body,omitempty"`
	ReadAt         *time.Time `json:"read_at,omitempty"`

	models.CommonTimestampsField
}

// Create 创建消息
func (m *Message) Create() {
	database.DB.Create(&m)
}

// Save 保存消息
func (m *Message) Save() (rowsAffected int64) {
	result := database.DB.Save(&m)
	return result.RowsAffected
}
