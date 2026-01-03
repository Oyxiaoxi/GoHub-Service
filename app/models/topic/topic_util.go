package topic

import (
	"GoHub-Service/pkg/app"
	"GoHub-Service/pkg/database"
	"GoHub-Service/pkg/paginator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Get(idstr string) (topic Topic) {
	database.DB.
		Select("id", "title", "body", "user_id", "category_id", "like_count", "favorite_count", "view_count", "created_at", "updated_at").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email", "avatar")
		}).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "description")
		}).
		Where("id", idstr).First(&topic)
	return
}

func GetBy(field, value string) (topic Topic) {
	database.DB.Where("? = ?", field, value).First(&topic)
	return
}

func All() (topics []Topic) {
	database.DB.Find(&topics)
	return
}

func IsExist(field, value string) bool {
	var count int64
	database.DB.Model(Topic{}).Where("? = ?", field, value).Count(&count)
	return count > 0
}

func Paginate(c *gin.Context, perPage int) (topics []Topic, paging paginator.Paging) {
	query := database.DB.Model(Topic{}).
		Select("id", "title", "body", "user_id", "category_id", "like_count", "favorite_count", "view_count", "created_at", "updated_at").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email", "avatar")
		}).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "description")
		}).
		Order("created_at DESC")

	paging = paginator.Paginate(
		c,
		query,
		&topics,
		app.V1URL(database.TableName(&Topic{})),
		perPage,
	)
	return
}

// BatchCreate 批量创建话题（使用事务和批量插入优化）
func BatchCreate(topics []Topic) error {
	if len(topics) == 0 {
		return nil
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 使用 CreateInBatches 批量插入，每批100条
		if err := tx.CreateInBatches(&topics, 100).Error; err != nil {
			return err
		}
		return nil
	})
}

// BatchDelete 批量删除话题（使用事务）
func BatchDelete(ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id IN ?", ids).Delete(&Topic{}).Error; err != nil {
			return err
		}
		return nil
	})
}
