package v1

import (
	"GoHub-Service/app/requests"
	"GoHub-Service/app/services"
	"GoHub-Service/pkg/auth"
	"GoHub-Service/pkg/config"
	apperrors "GoHub-Service/pkg/errors"
	"GoHub-Service/pkg/logger"
	"GoHub-Service/pkg/response"

	"github.com/gin-gonic/gin"
)

// MessagesController 私信接口
type MessagesController struct {
	service *services.MessageService
}

// NewMessagesController 创建控制器
func NewMessagesController() *MessagesController {
	return &MessagesController{service: services.NewMessageService()}
}

// Send 发送私信
func (ctrl *MessagesController) Send(c *gin.Context) {
	request := requests.MessageSendRequest{}
	if ok := requests.Validate(c, &request, requests.MessageSend); !ok {
		return
	}

	currentUserID := auth.CurrentUID(c)
	msg, err := ctrl.service.Send(currentUserID, request.ReceiverID, request.Body)
	if err != nil {
		handleMessageError(c, err, "发送私信失败")
		return
	}

	response.Created(c, msg)
}

// Conversation 获取会话消息
func (ctrl *MessagesController) Conversation(c *gin.Context) {
	request := requests.MessageConversationRequest{}
	if ok := requests.Validate(c, &request, requests.MessageConversation); !ok {
		return
	}

	currentUserID := auth.CurrentUID(c)
	perPage := config.GetInt("paging.perpage")
	list, paging, unread, err := ctrl.service.Conversation(c, currentUserID, request.UserID, perPage)
	if err != nil {
		handleMessageError(c, err, "获取会话失败")
		return
	}

	response.JSON(c, gin.H{
		"data":   list,
		"paging": paging,
		"unread": unread,
	})
}

// MarkRead 标记会话已读
func (ctrl *MessagesController) MarkRead(c *gin.Context) {
	request := requests.MessageMarkReadRequest{}
	if ok := requests.Validate(c, &request, requests.MessageMarkRead); !ok {
		return
	}

	currentUserID := auth.CurrentUID(c)
	if _, err := ctrl.service.MarkRead(currentUserID, request.UserID); err != nil {
		handleMessageError(c, err, "标记已读失败")
		return
	}

	response.Success(c)
}

// UnreadCount 获取未读数量
func (ctrl *MessagesController) UnreadCount(c *gin.Context) {
	currentUserID := auth.CurrentUID(c)
	count, err := ctrl.service.CountUnread(currentUserID)
	if err != nil {
		handleMessageError(c, err, "统计未读失败")
		return
	}
	response.JSON(c, gin.H{"unread": count})
}

func handleMessageError(c *gin.Context, err *apperrors.AppError, message string) {
	if err == nil {
		return
	}
	switch err.Type {
	case apperrors.ErrorTypeValidation:
		response.ValidationError(c, map[string][]string{"error": {err.Message}})
	case apperrors.ErrorTypeNotFound:
		response.Abort404(c)
	case apperrors.ErrorTypeAuthorization:
		response.Unauthorized(c, err.Message)
	default:
		logger.LogErrorWithContext(c, err, message)
		response.ApiError(c, 500, err.Code, err.Message)
	}
}
