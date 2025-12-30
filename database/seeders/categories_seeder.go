package seeders

import (
	"fmt"

	"GoHub-Service/app/models/category"
	"GoHub-Service/pkg/console"
	"GoHub-Service/pkg/logger"
	"GoHub-Service/pkg/seed"

	"gorm.io/gorm"
)

func init() {

	seed.Add("SeedCategoriesTable", func(db *gorm.DB) {
		cats := []category.Category{
			{Name: "Tech", Description: "后端/前端/移动/云原生"},
			{Name: "AI", Description: "LLM、CV、NLP 与应用"},
			{Name: "Career", Description: "职场、面试、成长"},
			{Name: "Life", Description: "生活方式、运动、旅行"},
			{Name: "Product", Description: "产品设计、体验优化"},
			{Name: "Ops", Description: "运维、SRE、可观测性"},
		}

		result := db.Table("categories").Create(&cats)
		if err := result.Error; err != nil {
			logger.LogIf(err)
			return
		}

		console.Success(fmt.Sprintf("Table [%v] %v rows seeded", result.Statement.Table, result.RowsAffected))
	})
}
