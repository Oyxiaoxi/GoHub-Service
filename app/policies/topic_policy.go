// Package policies 用户授权
package policies

import (
    "GoHub-Service/app/models/topic"
    "GoHub-Service/pkg/auth"

    "github.com/gin-gonic/gin"
)

// CanModifyTopic 判断用户是否有权限修改话题
func CanModifyTopic(c *gin.Context, _topic topic.Topic) bool {
    return auth.CurrentUID(c) == _topic.UserID
}
