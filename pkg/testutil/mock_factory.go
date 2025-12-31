// Package testutil 提供测试数据工厂
package testutil

import (
	"GoHub-Service/app/models/category"
	"GoHub-Service/app/models/comment"
	"GoHub-Service/app/models/topic"
	"GoHub-Service/app/models/user"
	"time"
)

// MockUserFactory 创建测试用户
func MockUserFactory(id, name, email string) *user.User {
	now := time.Now()
	return &user.User{
		ID:              id,
		Name:            name,
		Email:           email,
		Phone:           "13800138000",
		Password:        "$2a$10$abcdefghijklmnopqrstuvwxyz", // bcrypt hash
		Avatar:          "https://example.com/avatar.jpg",
		Introduction:    "测试用户简介",
		NotificationCount: 0,
		LastActiveAt:    &now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// MockCategoryFactory 创建测试分类
func MockCategoryFactory(id, name, description string) *category.Category {
	now := time.Now()
	return &category.Category{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// MockTopicFactory 创建测试话题
func MockTopicFactory(id, title, body, userID, categoryID string) *topic.Topic {
	now := time.Now()
	return &topic.Topic{
		ID:            id,
		Title:         title,
		Body:          body,
		UserID:        userID,
		CategoryID:    categoryID,
		ViewCount:     0,
		LikeCount:     0,
		CommentCount:  0,
		FavoriteCount: 0,
		LastCommentAt: &now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// MockCommentFactory 创建测试评论
func MockCommentFactory(id, content, topicID, userID, parentID string) *comment.Comment {
	now := time.Now()
	return &comment.Comment{
		ID:        id,
		Content:   content,
		TopicID:   topicID,
		UserID:    userID,
		ParentID:  parentID,
		LikeCount: 0,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// MockCategories 创建多个测试分类
func MockCategories() []category.Category {
	return []category.Category{
		*MockCategoryFactory("1", "技术", "技术讨论分类"),
		*MockCategoryFactory("2", "生活", "生活分享分类"),
		*MockCategoryFactory("3", "娱乐", "娱乐八卦分类"),
	}
}

// MockTopics 创建多个测试话题
func MockTopics() []topic.Topic {
	return []topic.Topic{
		*MockTopicFactory("1", "Go语言最佳实践", "分享Go语言开发经验...", "1", "1"),
		*MockTopicFactory("2", "如何提高代码质量", "代码质量提升技巧...", "1", "1"),
		*MockTopicFactory("3", "周末去哪玩", "推荐好玩的地方...", "2", "2"),
	}
}

// MockComments 创建多个测试评论
func MockComments() []comment.Comment {
	return []comment.Comment{
		*MockCommentFactory("1", "很有用的内容！", "1", "2", ""),
		*MockCommentFactory("2", "学到了很多", "1", "3", ""),
		*MockCommentFactory("3", "同意楼上", "1", "4", "1"),
	}
}

// MockTime 返回固定的测试时间
func MockTime() time.Time {
	return time.Date(2025, 12, 31, 12, 0, 0, 0, time.UTC)
}

// MockTimePtr 返回固定的测试时间指针
func MockTimePtr() *time.Time {
	t := MockTime()
	return &t
}
