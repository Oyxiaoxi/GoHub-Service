package migrations

import (
	"GoHub-Service/app/models"
	"GoHub-Service/pkg/migrate"
	"database/sql"

	"gorm.io/gorm"
)

func init() {

	type Link struct {
		models.BaseModel

		Name string `gorm:"type:varchar(255);not null"`
		URL  string `gorm:"type:varchar(255);default:null"`

		models.CommonTimestampsField
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		migrator.AutoMigrate(&Link{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		migrator.DropTable(&Link{})
	}

	migrate.Add("2022_01_13_145923_add_links_table", up, down)
}
