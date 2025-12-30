package seeders

import (
	"fmt"
	"math/rand"
	"time"

	"GoHub-Service/app/models"
	"GoHub-Service/app/models/comment"
	"GoHub-Service/pkg/console"
	"GoHub-Service/pkg/database"
	"GoHub-Service/pkg/logger"
	"GoHub-Service/pkg/seed"

	"github.com/bxcodec/faker/v3"
	"gorm.io/gorm"
)

func init() {

	seed.Add("SeedCommentsTable", func(db *gorm.DB) {
		rand.Seed(time.Now().UnixNano())
		var topicIDs []string
		var userIDs []string
		database.DB.Table("topics").Select("id").Find(&topicIDs)
		database.DB.Table("users").Select("id").Find(&userIDs)
		if len(topicIDs) == 0 || len(userIDs) == 0 {
			logger.LogIf(fmt.Errorf("seed comments requires topics and users"))
			return
		}

		var topLevel []comment.Comment
		total := 120
		for i := 0; i < total; i++ {
			topicID := topicIDs[rand.Intn(len(topicIDs))]
			author := userIDs[rand.Intn(len(userIDs))]
			created := time.Now().Add(-time.Duration(rand.Intn(60*24)) * time.Hour)
			topLevel = append(topLevel, comment.Comment{
				TopicID:   topicID,
				UserID:    author,
				Content:   faker.Sentence(),
				ParentID:  "0",
				LikeCount: int64(rand.Intn(40)),
				CommonTimestampsField: models.CommonTimestampsField{
					CreatedAt: created,
					UpdatedAt: created,
				},
			})
		}

		if err := db.Table("comments").Create(&topLevel).Error; err != nil {
			logger.LogIf(err)
			return
		}

		var parentIDs []string
		database.DB.Table("comments").Where("parent_id = ?", "0").Pluck("id", &parentIDs)
		if len(parentIDs) == 0 {
			console.Success("Comments seeded (no replies)")
			return
		}

		var replies []comment.Comment
		replyCount := 25
		for i := 0; i < replyCount; i++ {
			topicID := topicIDs[rand.Intn(len(topicIDs))]
			author := userIDs[rand.Intn(len(userIDs))]
			created := time.Now().Add(-time.Duration(rand.Intn(30*24)) * time.Hour)
			parentID := parentIDs[rand.Intn(len(parentIDs))]
			replies = append(replies, comment.Comment{
				TopicID:   topicID,
				UserID:    author,
				Content:   faker.Sentence(),
				ParentID:  parentID,
				LikeCount: int64(rand.Intn(20)),
				CommonTimestampsField: models.CommonTimestampsField{
					CreatedAt: created,
					UpdatedAt: created,
				},
			})
		}

		result := db.Table("comments").Create(&replies)
		if err := result.Error; err != nil {
			logger.LogIf(err)
			return
		}

		console.Success(fmt.Sprintf("Table [%v] %v rows seeded (including replies)", result.Statement.Table, len(topLevel)+len(replies)))
	})
}
