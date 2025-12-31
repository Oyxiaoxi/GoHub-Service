package routes

import (
	"GoHub-Service/app/http/controllers/admin"
	"GoHub-Service/app/http/middlewares"

	"github.com/gin-gonic/gin"
)

// RegisterAdminRoutes 注册管理后台路由
func RegisterAdminRoutes(r *gin.Engine) {
	// 管理后台路由组
	// 所有路由都需要认证，并且需要 admin 角色
	adminGroup := r.Group("/api/v1/admin")
	adminGroup.Use(middlewares.AuthJWT())
	adminGroup.Use(middlewares.RequireRole("admin"))
	{
		// 管理后台根路径
		adminGroup.GET("", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "GoHub 管理后台 API",
				"version": "v1",
				"endpoints": gin.H{
					"dashboard": "/api/v1/admin/dashboard/overview",
					"users":     "/api/v1/admin/users",
					"topics":    "/api/v1/admin/topics",
					"categories": "/api/v1/admin/categories",
				},
			})
		})

		// 仪表盘
		dashboardController := &admin.DashboardController{}
		adminGroup.GET("/dashboard/overview", dashboardController.Overview)
		adminGroup.GET("/dashboard/recent-users", dashboardController.RecentUsers)
		adminGroup.GET("/dashboard/recent-topics", dashboardController.RecentTopics)

		// 用户管理
		userController := &admin.UserController{}
		users := adminGroup.Group("/users")
		{
			users.GET("", userController.Index)           // 用户列表
			users.GET("/:id", userController.Show)        // 用户详情
			users.PUT("/:id", userController.Update)      // 更新用户
			users.DELETE("/:id", userController.Delete)   // 删除用户
			users.POST("/batch-delete", userController.BatchDelete) // 批量删除
			
			// 用户操作
			users.POST("/:id/ban", userController.Ban)              // 封禁用户
			users.POST("/:id/unban", userController.Unban)          // 解封用户
			users.POST("/:id/reset-password", userController.ResetPassword) // 重置密码
			users.POST("/:id/assign-role", userController.AssignRole)       // 分配角色
		}

		// 话题管理
		topicController := &admin.TopicController{}
		topics := adminGroup.Group("/topics")
		{
			topics.GET("", topicController.Index)         // 话题列表
			topics.GET("/:id", topicController.Show)      // 话题详情
			topics.PUT("/:id", topicController.Update)    // 更新话题
			topics.DELETE("/:id", topicController.Delete) // 删除话题
			topics.POST("/batch-delete", topicController.BatchDelete) // 批量删除
			
			// 话题操作
			topics.POST("/:id/pin", topicController.Pin)         // 置顶
			topics.POST("/:id/unpin", topicController.Unpin)     // 取消置顶
			topics.POST("/:id/approve", topicController.Approve) // 审核通过
			topics.POST("/:id/reject", topicController.Reject)   // 审核拒绝
		}

		// 分类管理
		categoryController := &admin.CategoryController{}
		categories := adminGroup.Group("/categories")
		{
			categories.GET("", categoryController.Index)         // 分类列表
			categories.GET("/:id", categoryController.Show)      // 分类详情
			categories.POST("", categoryController.Store)        // 创建分类
			categories.PUT("/:id", categoryController.Update)    // 更新分类
			categories.DELETE("/:id", categoryController.Delete) // 删除分类
			categories.POST("/sort", categoryController.Sort)    // 分类排序
		}

		// 角色管理
		roleController := &admin.RoleController{}
		roles := adminGroup.Group("/roles")
		{
			roles.GET("", roleController.Index)              // 获取角色列表
			roles.POST("", roleController.Store)             // 创建角色
			roles.GET("/:id", roleController.Show)           // 获取角色详情
			roles.PUT("/:id", roleController.Update)         // 更新角色
			roles.DELETE("/:id", roleController.Delete)      // 删除角色
			roles.GET("/:id/permissions", roleController.GetPermissions)       // 获取角色权限
			roles.POST("/:id/permissions", roleController.AssignPermissions)   // 分配权限
		}

		// 权限管理
		permissionController := &admin.PermissionController{}
		permissions := adminGroup.Group("/permissions")
		{
			permissions.GET("", permissionController.Index)         // 获取权限列表
			permissions.POST("", permissionController.Store)        // 创建权限
			permissions.GET("/:id", permissionController.Show)      // 获取权限详情
			permissions.PUT("/:id", permissionController.Update)    // 更新权限
			permissions.DELETE("/:id", permissionController.Delete) // 删除权限
		}
	}

	// 版主路由组（moderator 角色）
	moderatorGroup := r.Group("/api/v1/moderator")
	moderatorGroup.Use(middlewares.AuthJWT())
	moderatorGroup.Use(middlewares.RequireRole("moderator"))
	{
		topicController := &admin.TopicController{}
		
		// 版主可以审核话题
		moderatorGroup.GET("/topics", topicController.Index)
		moderatorGroup.POST("/topics/:id/approve", topicController.Approve)
		moderatorGroup.POST("/topics/:id/reject", topicController.Reject)
		moderatorGroup.DELETE("/topics/:id", topicController.Delete)
	}
}
