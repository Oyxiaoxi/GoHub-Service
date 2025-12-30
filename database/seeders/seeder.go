// Package seeders 存放数据填充文件
package seeders

import "GoHub-Service/pkg/seed"

func Initialize() {

	// 触发加载本目录下其他文件中的 init 方法

	// 指定优先于同目录下的其他文件运行
	seed.SetRunOrder([]string{
		// 顺序需满足外键依赖：用户 -> 角色/权限 -> 分类 -> 话题 -> 评论 -> 互动/私信 -> 友情链接
		"SeedUsersTable",
		"SeedRolesAndPermissions", // 角色和权限初始化
		"SeedCategoriesTable",
		"SeedTopicsTable",
		"SeedCommentsTable",
		"SeedInteractions",
		"SeedMessages",
		"SeedLinksTable",
	})
}
