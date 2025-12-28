package migrations

import (
	"database/sql"
	"GoHub-Service/pkg/migrate"

	"gorm.io/gorm"
)

func init() {
	up := func(migrator gorm.Migrator, DB *sql.DB) {
		// 为 topics 表添加性能索引
		DB.Exec("ALTER TABLE topics ADD INDEX idx_topics_created_at (created_at);")
		DB.Exec("ALTER TABLE topics ADD INDEX idx_topics_updated_at (updated_at);")
		
		// 为 users 表添加性能索引
		DB.Exec("ALTER TABLE users ADD INDEX idx_users_phone (phone);")
		DB.Exec("ALTER TABLE users ADD INDEX idx_users_email (email);")
		DB.Exec("ALTER TABLE users ADD INDEX idx_users_created_at (created_at);")
		
		// 为 categories 表添加性能索引
		DB.Exec("ALTER TABLE categories ADD INDEX idx_categories_created_at (created_at);")
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		// 删除索引
		DB.Exec("ALTER TABLE topics DROP INDEX idx_topics_created_at;")
		DB.Exec("ALTER TABLE topics DROP INDEX idx_topics_updated_at;")
		DB.Exec("ALTER TABLE users DROP INDEX idx_users_phone;")
		DB.Exec("ALTER TABLE users DROP INDEX idx_users_email;")
		DB.Exec("ALTER TABLE users DROP INDEX idx_users_created_at;")
		DB.Exec("ALTER TABLE categories DROP INDEX idx_categories_created_at;")
	}

	migrate.Add("2025_12_29_004018_add_performance_indexes", up, down)
}
