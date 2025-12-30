package services

import (
	"GoHub-Service/app/repositories"
	"GoHub-Service/pkg/config"

	"github.com/gin-gonic/gin"
)

// SearchService 搜索业务逻辑
type SearchService interface {
	SearchTopics(c *gin.Context, keyword string) (interface{}, error)
	SearchUsers(c *gin.Context, keyword string) (interface{}, error)
}

type searchService struct {
	repo repositories.SearchRepository
}

// NewSearchService 创建实例
func NewSearchService(repo repositories.SearchRepository) SearchService {
	return &searchService{repo: repo}
}

func (s *searchService) SearchTopics(c *gin.Context, keyword string) (interface{}, error) {
	perPage := config.GetInt("paging.perpage")
	topics, paging, err := s.repo.SearchTopics(c, keyword, perPage)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"data":   topics,
		"paging": paging,
	}, nil
}

func (s *searchService) SearchUsers(c *gin.Context, keyword string) (interface{}, error) {
	perPage := config.GetInt("paging.perpage")
	users, paging, err := s.repo.SearchUsers(c, keyword, perPage)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"data":   users,
		"paging": paging,
	}, nil
}
