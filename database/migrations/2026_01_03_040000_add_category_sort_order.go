package migrations

import (
	"database/sql"

	"GoHub-Service/pkg/migrate"

	"gorm.io/gorm"
)

func init() {
	type Category struct {
		SortOrder int `gorm:"type:int;default:0;index;comment:排序顺序"`
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.AutoMigrate(&Category{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.DropColumn(&Category{}, "sort_order")
	}

	migrate.Add("2026_01_03_040000_add_category_sort_order", up, down)
}
