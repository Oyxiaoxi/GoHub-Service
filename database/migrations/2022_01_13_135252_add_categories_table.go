package migrations

import (
	"GoHub-Service/app/models"
	"GoHub-Service/pkg/migrate"
	"database/sql"

	"gorm.io/gorm"
)

func init() {

	type Category struct {
		models.BaseModel

		Name        string `gorm:"type:varchar(255);not null;index"`
		Description string `gorm:"type:varchar(255);default:null"`

		models.CommonTimestampsField
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		migrator.AutoMigrate(&Category{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		migrator.DropTable(&Category{})
	}

	migrate.Add("2022_01_13_135252_add_categories_table", up, down)
}
