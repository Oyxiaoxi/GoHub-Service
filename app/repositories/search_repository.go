// Package repositories 搜索数据访问层
package repositories

import (
	"GoHub-Service/app/models/topic"
	"GoHub-Service/app/models/user"
	"GoHub-Service/pkg/database"
	"GoHub-Service/pkg/paginator"

	"github.com/gin-gonic/gin"
)

// SearchRepository 搜索仓储接口
type SearchRepository interface {
	SearchTopics(c *gin.Context, keyword string, perPage int) ([]topic.Topic, *paginator.Paging, error)
	SearchUsers(c *gin.Context, keyword string, perPage int) ([]user.User, *paginator.Paging, error)
}

type searchRepository struct{}

// NewSearchRepository 创建实例
func NewSearchRepository() SearchRepository {
	return &searchRepository{}
}

func (r *searchRepository) SearchTopics(c *gin.Context, keyword string, perPage int) ([]topic.Topic, *paginator.Paging, error) {
	var topics []topic.Topic
	like := "%" + keyword + "%"
	query := database.DB.Model(&topic.Topic{}).
		Where("title LIKE ? OR body LIKE ?", like, like).
		Preload("User").
		Preload("Category").
		Order("created_at DESC")

	paging := paginator.Paginate(c, query, &topics, "/api/v1/search/topics", perPage)
	return topics, &paging, nil
}

func (r *searchRepository) SearchUsers(c *gin.Context, keyword string, perPage int) ([]user.User, *paginator.Paging, error) {
	var users []user.User
	like := "%" + keyword + "%"
	query := database.DB.Model(&user.User{}).
		Where("name LIKE ? OR city LIKE ? OR introduction LIKE ?", like, like, like).
		Order("created_at DESC")

	paging := paginator.Paginate(c, query, &users, "/api/v1/search/users", perPage)
	return users, &paging, nil
}
