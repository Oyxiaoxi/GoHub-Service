package services

import (
	"strings"

	"GoHub-Service/app/models/message"
	"GoHub-Service/app/repositories"
	apperrors "GoHub-Service/pkg/errors"
	"GoHub-Service/pkg/paginator"

	"github.com/gin-gonic/gin"
)

// MessageService 私信业务
type MessageService struct {
	repo     repositories.MessageRepository
	userRepo repositories.UserRepository
	notifSvc *NotificationService
}

// NewMessageService 创建实例
func NewMessageService() *MessageService {
	return &MessageService{
		repo:     repositories.NewMessageRepository(),
		userRepo: repositories.NewUserRepository(),
		notifSvc: NewNotificationService(),
	}
}

// Send 发送私信
func (s *MessageService) Send(senderID, receiverID, body string) (*message.Message, *apperrors.AppError) {
	if senderID == "" {
		return nil, apperrors.AuthorizationError("未登录")
	}
	if senderID == receiverID {
		return nil, apperrors.ValidationError("不能给自己发送消息", map[string]interface{}{"receiver_id": receiverID})
	}

	// 确保收件人存在
	if _, err := s.userRepo.GetByID(receiverID); err != nil {
		return nil, apperrors.WrapError(err, "获取收件人失败")
	}

	convID := buildConversationID(senderID, receiverID)
	msg := &message.Message{
		ConversationID: convID,
		SenderID:       senderID,
		ReceiverID:     receiverID,
		Body:           strings.TrimSpace(body),
	}

	if err := s.repo.Create(msg); err != nil {
		return nil, apperrors.DatabaseError("创建私信", err)
	}

	if s.notifSvc != nil {
		_ = s.notifSvc.Notify(receiverID, senderID, "direct_message", map[string]interface{}{"message_id": msg.GetStringID(), "sender_id": senderID})
	}

	return msg, nil
}

// Conversation 获取双方的会话消息
func (s *MessageService) Conversation(c *gin.Context, currentUserID, partnerID string, perPage int) ([]message.Message, *paginator.Paging, int64, *apperrors.AppError) {
	if currentUserID == partnerID {
		return nil, nil, 0, apperrors.ValidationError("不能查看与自己的会话", map[string]interface{}{"user_id": partnerID})
	}

	if _, err := s.userRepo.GetByID(partnerID); err != nil {
		return nil, nil, 0, apperrors.WrapError(err, "获取会话用户失败")
	}

	convID := buildConversationID(currentUserID, partnerID)
	list, paging, err := s.repo.ListConversation(c, convID, currentUserID, perPage)
	if err != nil {
		return nil, nil, 0, apperrors.DatabaseError("获取会话消息", err)
	}

	// 标记为已读
	if _, err := s.repo.MarkConversationRead(convID, currentUserID); err != nil {
		return nil, nil, 0, apperrors.DatabaseError("标记已读", err)
	}

	unread, err := s.repo.CountUnread(currentUserID)
	if err != nil {
		return nil, nil, 0, apperrors.DatabaseError("统计未读消息", err)
	}

	return list, paging, unread, nil
}

// MarkRead 标记会话为已读
func (s *MessageService) MarkRead(currentUserID, partnerID string) (int64, *apperrors.AppError) {
	if _, err := s.userRepo.GetByID(partnerID); err != nil {
		return 0, apperrors.WrapError(err, "获取会话用户失败")
	}

	convID := buildConversationID(currentUserID, partnerID)
	affected, err := s.repo.MarkConversationRead(convID, currentUserID)
	if err != nil {
		return 0, apperrors.DatabaseError("标记已读", err)
	}
	return affected, nil
}

// CountUnread 统计未读
func (s *MessageService) CountUnread(currentUserID string) (int64, *apperrors.AppError) {
	unread, err := s.repo.CountUnread(currentUserID)
	if err != nil {
		return 0, apperrors.DatabaseError("统计未读消息", err)
	}
	return unread, nil
}

func buildConversationID(a, b string) string {
	if a < b {
		return a + ":" + b
	}
	return b + ":" + a
}
