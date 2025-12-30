package seeders

import (
	"fmt"
	"GoHub-Service/database/factories"
	"GoHub-Service/pkg/console"
	"GoHub-Service/pkg/logger"
	"GoHub-Service/pkg/seed"

	"gorm.io/gorm"
)

func init() {

	seed.Add("SeedCommentsTable", func(db *gorm.DB) {

		comments := factories.MakeComments(50)

		result := db.Table("comments").Create(&comments)

		if err := result.Error; err != nil {
			logger.LogIf(err)
			return
		}

		console.Success(fmt.Sprintf("Table [%v] %v rows seeded", result.Statement.Table, result.RowsAffected))
	})
}
