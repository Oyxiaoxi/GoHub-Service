package services_test

import (
	"GoHub-Service/app/services"
	"testing"
	"time"

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
	title := "更新话题"
	body := "更新内容"
	categoryID := "2"
	
	dto := services.TopicUpdateDTO{
		Title:      &title,
		Body:       &body,
		CategoryID: &categoryID,
	}

	assert.NotNil(t, dto.Title)
	assert.Equal(t, "更新话题", *dto.Title)
	assert.NotNil(t, dto.Body)
	assert.Equal(t, "更新内容", *dto.Body)
	assert.NotNil(t, dto.CategoryID)
	assert.Equal(t, "2", *dto.CategoryID)
}

// TestTopicResponseDTO 测试话题响应DTO
func TestTopicResponseDTO(t *testing.T) {
	now := time.Now()
	dto := services.TopicResponseDTO{
		ID:         "1",
		Title:      "测试话题",
		Body:       "测试内容",
		CategoryID: "1",
		UserID:     "1",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	assert.Equal(t, "1", dto.ID)
	assert.Equal(t, "测试话题", dto.Title)
	assert.Equal(t, "测试内容", dto.Body)
	assert.NotEmpty(t, dto.CreatedAt)
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
	name := "更新分类"
	description := "更新描述"
	
	dto := services.CategoryUpdateDTO{
		Name:        &name,
		Description: &description,
	}

	assert.NotNil(t, dto.Name)
	assert.Equal(t, "更新分类", *dto.Name)
	assert.NotNil(t, dto.Description)
	assert.Equal(t, "更新描述", *dto.Description)
}

// TestCategoryResponseDTO 测试分类响应DTO
func TestCategoryResponseDTO(t *testing.T) {
	now := time.Now()
	dto := services.CategoryResponseDTO{
		ID:          "1",
		Name:        "测试分类",
		Description: "测试描述",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	assert.Equal(t, "1", dto.ID)
	assert.Equal(t, "测试分类", dto.Name)
	assert.Equal(t, "测试描述", dto.Description)
	assert.NotEmpty(t, dto.CreatedAt)
}

// TestUserCreateDTO 测试创建用户DTO
func TestUserCreateDTO(t *testing.T) {
	dto := services.UserCreateDTO{
		Name:     "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Phone:    "13800138000",
	}

	assert.Equal(t, "testuser", dto.Name)
	assert.Equal(t, "test@example.com", dto.Email)
	assert.Equal(t, "password123", dto.Password)
	assert.Equal(t, "13800138000", dto.Phone)
}

// TestUserUpdateDTO 测试更新用户DTO
func TestUserUpdateDTO(t *testing.T) {
	name := "updateduser"
	email := "updated@example.com"
	phone := "13900139000"
	
	dto := services.UserUpdateDTO{
		Name:  &name,
		Email: &email,
		Phone: &phone,
	}

	assert.NotNil(t, dto.Name)
	assert.Equal(t, "updateduser", *dto.Name)
	assert.NotNil(t, dto.Email)
	assert.Equal(t, "updated@example.com", *dto.Email)
	assert.NotNil(t, dto.Phone)
	assert.Equal(t, "13900139000", *dto.Phone)
}

// TestUserResponseDTO 测试用户响应DTO
func TestUserResponseDTO(t *testing.T) {
	now := time.Now()
	dto := services.UserResponseDTO{
		ID:        "1",
		Name:      "testuser",
		Email:     "test@example.com",
		Phone:     "13800138000",
		CreatedAt: now,
		UpdatedAt: now,
	}

	assert.Equal(t, "1", dto.ID)
	assert.Equal(t, "testuser", dto.Name)
	assert.Equal(t, "test@example.com", dto.Email)
	assert.NotEmpty(t, dto.CreatedAt)
}

// TestLinkResponseDTO 测试友情链接响应DTO
func TestLinkResponseDTO(t *testing.T) {
	now := time.Now()
	dto := services.LinkResponseDTO{
		ID:        "1",
		Name:      "测试链接",
		URL:       "https://example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}

	assert.Equal(t, "1", dto.ID)
	assert.Equal(t, "测试链接", dto.Name)
	assert.Equal(t, "https://example.com", dto.URL)
	assert.NotEmpty(t, dto.CreatedAt)
}
