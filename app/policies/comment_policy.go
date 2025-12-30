// Package policies 评论授权
package policies

import (
	"GoHub-Service/app/models/comment"
	"GoHub-Service/pkg/auth"

	"github.com/gin-gonic/gin"
)

// CanModifyComment 判断用户是否有权限修改评论
func CanModifyComment(c *gin.Context, _comment comment.Comment) bool {
	return auth.CurrentUID(c) == _comment.UserID
}
