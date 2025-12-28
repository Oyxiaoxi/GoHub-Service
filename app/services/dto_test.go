package services_test

import (
	"GoHub-Service/app/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTopicCreateDTO 测试创建话题DTO
func TestTopicCreateDTO(t *testing.T) {
	dto := services.TopicCreateDTO{
		Title:      "测试话题",
		Body:       "测试内容",
		CategoryID: "1",
		UserID:     "1",
	}

	assert.Equal(t, "测试话题", dto.Title)
	assert.Equal(t, "测试内容", dto.Body)
	assert.Equal(t, "1", dto.CategoryID)
	assert.Equal(t, "1", dto.UserID)
}

// TestTopicUpdateDTO 测试更新话题DTO
func TestTopicUpdateDTO(t *testing.T) {
	dto := services.TopicUpdateDTO{
		Title:      "更新话题",
		Body:       "更新内容",
		CategoryID: "2",
	}

	assert.Equal(t, "更新话题", dto.Title)
	assert.Equal(t, "更新内容", dto.Body)
	assert.Equal(t, "2", dto.CategoryID)
}

// TestCategoryCreateDTO 测试创建分类DTO
func TestCategoryCreateDTO(t *testing.T) {
	dto := services.CategoryCreateDTO{
		Name:        "测试分类",
		Description: "测试描述",
	}

	assert.Equal(t, "测试分类", dto.Name)
	assert.Equal(t, "测试描述", dto.Description)
}

// TestCategoryUpdateDTO 测试更新分类DTO
func TestCategoryUpdateDTO(t *testing.T) {
	dto := services.CategoryUpdateDTO{
		Name:        "更新分类",
		Description: "更新描述",
	}

	assert.Equal(t, "更新分类", dto.Name)
	assert.Equal(t, "更新描述", dto.Description)
}
