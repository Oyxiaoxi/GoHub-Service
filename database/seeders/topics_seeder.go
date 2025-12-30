package seeders

import (
	"fmt"
	"math/rand"
	"time"

	"GoHub-Service/app/models"
	"GoHub-Service/app/models/topic"
	"GoHub-Service/pkg/console"
	"GoHub-Service/pkg/database"
	"GoHub-Service/pkg/logger"
	"GoHub-Service/pkg/seed"

	"github.com/bxcodec/faker/v3"
	"gorm.io/gorm"
)

func init() {

	seed.Add("SeedTopicsTable", func(db *gorm.DB) {
		rand.Seed(time.Now().UnixNano())

		var userIDs []string
		var categoryIDs []string
		database.DB.Table("users").Select("id").Find(&userIDs)
		database.DB.Table("categories").Select("id").Find(&categoryIDs)
		if len(userIDs) == 0 || len(categoryIDs) == 0 {
			logger.LogIf(fmt.Errorf("seed topics requires users and categories"))
			return
		}

		count := 40
		var topics []topic.Topic
		for i := 0; i < count; i++ {
			author := userIDs[rand.Intn(len(userIDs))]
			cat := categoryIDs[rand.Intn(len(categoryIDs))]
			created := time.Now().Add(-time.Duration(rand.Intn(90*24)) * time.Hour)
			title := faker.Sentence()
			body := faker.Paragraph()
			topics = append(topics, topic.Topic{
				Title:         title,
				Body:          body,
				CategoryID:    cat,
				UserID:        author,
				ViewCount:     int64(rand.Intn(2000)),
				LikeCount:     int64(rand.Intn(200)),
				FavoriteCount: int64(rand.Intn(150)),
				CommonTimestampsField: models.CommonTimestampsField{
					CreatedAt: created,
					UpdatedAt: created,
				},
			})
		}

		result := db.Table("topics").Create(&topics)
		if err := result.Error; err != nil {
			logger.LogIf(err)
			return
		}

		console.Success(fmt.Sprintf("Table [%v] %v rows seeded", result.Statement.Table, result.RowsAffected))
	})
}
