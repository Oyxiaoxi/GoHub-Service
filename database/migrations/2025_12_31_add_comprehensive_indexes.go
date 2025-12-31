package migrations

import (
	"database/sql"
	"GoHub-Service/pkg/migrate"

	"gorm.io/gorm"
)

func init() {

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		// ============================================================
		// Topics 表索引优化
		// ============================================================
		
		// 复合索引：按用户+时间查询话题
		DB.Exec("CREATE INDEX idx_topics_user_created ON topics(user_id, created_at DESC)")
		
		// 复合索引：按分类+点赞数查询热门话题
		DB.Exec("CREATE INDEX idx_topics_category_likes ON topics(category_id, like_count DESC)")
		
		// ============================================================
		// Comments 表索引优化
		// ============================================================
		
		// 复合索引：查询话题的顶级评论并按时间排序
		DB.Exec("CREATE INDEX idx_comments_topic_parent_created ON comments(topic_id, parent_id, created_at DESC)")
		
		// 复合索引：查询用户评论按时间排序
		DB.Exec("CREATE INDEX idx_comments_user_created ON comments(user_id, created_at DESC)")
		
		// 复合索引：查询父评论的回复
		DB.Exec("CREATE INDEX idx_comments_parent_created ON comments(parent_id, created_at ASC)")
		
		// ============================================================
		// User_Roles 表索引优化（RBAC权限检查）
		// ============================================================
		
		// 单列索引：优化用户权限查询
		DB.Exec("CREATE INDEX idx_user_roles_user_id ON user_roles(user_id)")
		
		// 单列索引：优化角色用户列表查询
		DB.Exec("CREATE INDEX idx_user_roles_role_id ON user_roles(role_id)")
		
		// ============================================================
		// Role_Permissions 表索引优化（RBAC权限检查）
		// ============================================================
		
		// 单列索引：优化角色权限查询
		DB.Exec("CREATE INDEX idx_role_permissions_role ON role_permissions(role_id)")
		DB.Exec("CREATE INDEX idx_role_permissions_permission ON role_permissions(permission_id)")
		
		// ============================================================
		// Topic_Likes 表索引优化（点赞查询）
		// ============================================================
		
		// 复合索引：用户点赞列表按时间排序
		DB.Exec("CREATE INDEX idx_topic_likes_user_created ON topic_likes(user_id, created_at DESC)")
		
		// 复合索引：话题点赞列表按时间排序
		DB.Exec("CREATE INDEX idx_topic_likes_topic_created ON topic_likes(topic_id, created_at DESC)")
		
		// ============================================================
		// Topic_Favorites 表索引优化（收藏查询）
		// ============================================================
		
		// 复合索引：用户收藏列表按时间排序
		DB.Exec("CREATE INDEX idx_topic_favorites_user_created ON topic_favorites(user_id, created_at DESC)")
		
		// 复合索引：话题收藏列表按时间排序
		DB.Exec("CREATE INDEX idx_topic_favorites_topic_created ON topic_favorites(topic_id, created_at DESC)")
		
		// ============================================================
		// User_Follows 表索引优化（关注关系查询）
		// ============================================================
		
		// 复合索引：用户的关注列表按时间排序
		DB.Exec("CREATE INDEX idx_user_follows_follower_created ON user_follows(follower_id, created_at DESC)")
		
		// 复合索引：用户的粉丝列表按时间排序
		DB.Exec("CREATE INDEX idx_user_follows_followee_created ON user_follows(followee_id, created_at DESC)")
		
		// ============================================================
		// Users 表索引优化
		// ============================================================
		
		// 单列索引：按积分排序（已存在，跳过）
		// DB.Exec("CREATE INDEX idx_users_points ON users(points DESC)")
		
		// 单列索引：按最后活跃时间排序
		DB.Exec("CREATE INDEX idx_users_last_active ON users(last_active_at DESC)")
		
		// ============================================================
		// Categories 表索引优化
		// ============================================================
		
		// 单列索引：按名称查询分类（已存在，跳过）
		// DB.Exec("CREATE INDEX idx_categories_name ON categories(name)")
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		// 删除所有新增的索引
		
		// Topics 表
		DB.Exec("DROP INDEX idx_topics_user_created ON topics")
		DB.Exec("DROP INDEX idx_topics_category_likes ON topics")
		
		// Comments 表
		DB.Exec("DROP INDEX idx_comments_topic_parent_created ON comments")
		DB.Exec("DROP INDEX idx_comments_user_created ON comments")
		DB.Exec("DROP INDEX idx_comments_parent_created ON comments")
		
		// User_Roles 表
		DB.Exec("DROP INDEX idx_user_roles_user_id ON user_roles")
		DB.Exec("DROP INDEX idx_user_roles_role_id ON user_roles")
		
		// Role_Permissions 表
		DB.Exec("DROP INDEX idx_role_permissions_role ON role_permissions")
		DB.Exec("DROP INDEX idx_role_permissions_permission ON role_permissions")
		
		// Topic_Likes 表
		DB.Exec("DROP INDEX idx_topic_likes_user_created ON topic_likes")
		DB.Exec("DROP INDEX idx_topic_likes_topic_created ON topic_likes")
		
		// Topic_Favorites 表
		DB.Exec("DROP INDEX idx_topic_favorites_user_created ON topic_favorites")
		DB.Exec("DROP INDEX idx_topic_favorites_topic_created ON topic_favorites")
		
		// User_Follows 表
		DB.Exec("DROP INDEX idx_user_follows_follower_created ON user_follows")
		DB.Exec("DROP INDEX idx_user_follows_followee_created ON user_follows")
		
		// Users 表
		// DB.Exec("DROP INDEX idx_users_points ON users")
		DB.Exec("DROP INDEX idx_users_last_active ON users")
		
		// Categories 表
		// DB.Exec("DROP INDEX idx_categories_name ON categories")
	}

	migrate.Add("2025_12_31_add_comprehensive_indexes", up, down)
}
