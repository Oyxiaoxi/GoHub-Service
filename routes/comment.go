// Package routes 评论相关路由
package routes

import (
	"GoHub-Service/app/http/controllers/api/v1"
	"GoHub-Service/app/http/middlewares"
	"github.com/gin-gonic/gin"
)

// RegisterCommentRoutes 注册评论相关路由
func RegisterCommentRoutes(rg *gin.RouterGroup, commentsCtrl *v1.CommentsController) {
	// 评论基础路由
	commentsGroup := rg.Group("/comments")
	{
		commentsGroup.GET("", commentsCtrl.Index)
		commentsGroup.GET("/:id", commentsCtrl.Show)
		// 创建和更新评论应用内容安全检查
		commentsGroup.POST("", 
			middlewares.AuthJWT(), 
			middlewares.SensitiveWordFilter(),
			commentsCtrl.Store,
		)
		commentsGroup.PUT("/:id", 
			middlewares.AuthJWT(), 
			middlewares.SensitiveWordFilter(),
			commentsCtrl.Update,
		)
		commentsGroup.DELETE("/:id", middlewares.AuthJWT(), commentsCtrl.Delete)
		
		// 评论点赞
		commentsGroup.POST("/:id/like", middlewares.AuthJWT(), commentsCtrl.Like)
		commentsGroup.POST("/:id/unlike", middlewares.AuthJWT(), commentsCtrl.Unlike)
		
		// 获取评论的回复
		commentsGroup.GET("/:id/replies", commentsCtrl.ListReplies)
	}

	// 话题的评论路由
	topicsGroup := rg.Group("/topics")
	{
		topicsGroup.GET("/:id/comments", commentsCtrl.ListByTopicID)
	}

	// 用户的评论路由
	usersGroup := rg.Group("/users")
	{
		usersGroup.GET("/:id/comments", commentsCtrl.ListByUserID)
	}
}
