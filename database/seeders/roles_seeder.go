package seeders

import (
	"GoHub-Service/app/models/permission"
	"GoHub-Service/app/models/role"
	"GoHub-Service/pkg/console"
	"GoHub-Service/pkg/database"
	"GoHub-Service/pkg/logger"
	"GoHub-Service/pkg/seed"
	"fmt"

	"gorm.io/gorm"
)

func init() {

	seed.Add("SeedRolesAndPermissions", func(db *gorm.DB) {
		// 创建角色
		roles := []role.Role{
			{
				Name:        "admin",
				DisplayName: "管理员",
				Description: "系统管理员，拥有所有权限",
			},
			{
				Name:        "moderator",
				DisplayName: "版主",
				Description: "内容审核员，可管理话题和评论",
			},
			{
				Name:        "user",
				DisplayName: "普通用户",
				Description: "普通注册用户",
			},
		}

		for _, r := range roles {
			if err := db.Create(&r).Error; err != nil {
				logger.LogIf(err)
				continue
			}
		}

		// 创建权限
		permissions := []permission.Permission{
			// 用户管理权限
			{Name: "users.view", DisplayName: "查看用户", Group: "users"},
			{Name: "users.edit", DisplayName: "编辑用户", Group: "users"},
			{Name: "users.delete", DisplayName: "删除用户", Group: "users"},
			{Name: "users.ban", DisplayName: "封禁用户", Group: "users"},

			// 话题管理权限
			{Name: "topics.view", DisplayName: "查看话题", Group: "topics"},
			{Name: "topics.create", DisplayName: "创建话题", Group: "topics"},
			{Name: "topics.edit", DisplayName: "编辑话题", Group: "topics"},
			{Name: "topics.delete", DisplayName: "删除话题", Group: "topics"},
			{Name: "topics.pin", DisplayName: "置顶话题", Group: "topics"},
			{Name: "topics.featured", DisplayName: "精华话题", Group: "topics"},

			// 评论管理权限
			{Name: "comments.view", DisplayName: "查看评论", Group: "comments"},
			{Name: "comments.create", DisplayName: "创建评论", Group: "comments"},
			{Name: "comments.edit", DisplayName: "编辑评论", Group: "comments"},
			{Name: "comments.delete", DisplayName: "删除评论", Group: "comments"},

			// 分类管理权限
			{Name: "categories.view", DisplayName: "查看分类", Group: "categories"},
			{Name: "categories.manage", DisplayName: "管理分类", Group: "categories"},

			// 系统管理权限
			{Name: "admin.settings", DisplayName: "系统设置", Group: "admin"},
			{Name: "admin.stats", DisplayName: "数据统计", Group: "admin"},
		}

		for _, p := range permissions {
			if err := db.Create(&p).Error; err != nil {
				logger.LogIf(err)
				continue
			}
		}

		// 为角色分配权限
		// 1. 管理员 - 拥有所有权限
		adminRole := role.GetByName("admin")
		if adminRole.ID > 0 {
			var allPermissions []permission.Permission
			database.DB.Find(&allPermissions)
			for _, perm := range allPermissions {
				db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", adminRole.ID, perm.ID)
			}
		}

		// 2. 版主 - 话题和评论管理权限
		moderatorRole := role.GetByName("moderator")
		if moderatorRole.ID > 0 {
			modPermNames := []string{
				"topics.view", "topics.edit", "topics.delete", "topics.pin", "topics.featured",
				"comments.view", "comments.delete",
				"categories.view",
			}
			for _, permName := range modPermNames {
				perm := permission.GetByName(permName)
				if perm.ID > 0 {
					db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", moderatorRole.ID, perm.ID)
				}
			}
		}

		// 3. 普通用户 - 基础权限
		userRole := role.GetByName("user")
		if userRole.ID > 0 {
			userPermNames := []string{
				"topics.view", "topics.create", "topics.edit",
				"comments.view", "comments.create", "comments.edit",
				"categories.view",
			}
			for _, permName := range userPermNames {
				perm := permission.GetByName(permName)
				if perm.ID > 0 {
					db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", userRole.ID, perm.ID)
				}
			}
		}

		// 给第一个用户分配管理员角色（如果存在）
		var firstUserID uint64
		if err := database.DB.Table("users").Select("id").Order("id ASC").Limit(1).Scan(&firstUserID).Error; err == nil && firstUserID > 0 {
			db.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", firstUserID, adminRole.ID)
			console.Success(fmt.Sprintf("已将用户 ID=%d 设为管理员", firstUserID))
		}

		console.Success("角色和权限初始化完成")
	})
}
