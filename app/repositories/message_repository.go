// Package repositories 私信数据访问层
package repositories

import (
	"time"

	"GoHub-Service/app/models/message"
	"GoHub-Service/pkg/database"
	"GoHub-Service/pkg/paginator"

	"github.com/gin-gonic/gin"
)

// MessageRepository 私信仓储接口
type MessageRepository interface {
	Create(msg *message.Message) error
	ListConversation(c *gin.Context, conversationID string, participantID string, perPage int) ([]message.Message, *paginator.Paging, error)
	MarkConversationRead(conversationID, receiverID string) (int64, error)
	CountUnread(receiverID string) (int64, error)
}

type messageRepository struct{}

// NewMessageRepository 创建实例
func NewMessageRepository() MessageRepository {
	return &messageRepository{}
}

func (r *messageRepository) Create(msg *message.Message) error {
	msg.Create()
	if msg.ID == 0 {
		return ErrCreateFailed
	}
	return nil
}

func (r *messageRepository) ListConversation(c *gin.Context, conversationID string, participantID string, perPage int) ([]message.Message, *paginator.Paging, error) {
	var messages []message.Message
	query := database.DB.Model(&message.Message{}).
		Where("conversation_id = ? AND (sender_id = ? OR receiver_id = ?)", conversationID, participantID, participantID).
		Order("created_at ASC")

	paging := paginator.Paginate(c, query, &messages, "/api/v1/messages", perPage)
	return messages, &paging, nil
}

func (r *messageRepository) MarkConversationRead(conversationID, receiverID string) (int64, error) {
	now := time.Now()
	result := database.DB.Model(&message.Message{}).
		Where("conversation_id = ? AND receiver_id = ? AND read_at IS NULL", conversationID, receiverID).
		Updates(map[string]interface{}{"read_at": &now})
	return result.RowsAffected, result.Error
}

func (r *messageRepository) CountUnread(receiverID string) (int64, error) {
	var count int64
	if err := database.DB.Model(&message.Message{}).
		Where("receiver_id = ? AND read_at IS NULL", receiverID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
