// Package testutil 提供完整的测试数据工厂
package testutil

import (
	"GoHub-Service/app/models/category"
	"GoHub-Service/app/models/comment"
	"GoHub-Service/app/models/link"
	"GoHub-Service/app/models/message"
	"GoHub-Service/app/models/notification"
	"GoHub-Service/app/models/permission"
	"GoHub-Service/app/models/role"
	"GoHub-Service/app/models/topic"
	"GoHub-Service/app/models/user"
	"time"
)

// ============== User Factories ==============

// MockUserFactory 创建测试用户
func MockUserFactory(id, name, email string) *user.User {
	now := MockTime()
	return &user.User{
		ID:        id,
		Name:      name,
		Email:     email,
		Phone:     "13800138000",
		Avatar:    "/uploads/default_avatar.png",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// MockUsers 批量创建测试用户
func MockUsers(count int) []user.User {
	users := make([]user.User, count)
	for i := 0; i < count; i++ {
		users[i] = *MockUserFactory(
			string(rune(i+1)),
			"测试用户"+string(rune(i+1)),
			"user"+string(rune(i+1))+"@test.com",
		)
	}
	return users
}

// ============== Category Factories ==============

// MockCategoryFactory 创建测试分类
func MockCategoryFactory(id, name, description string) *category.Category {
	now := MockTime()
	return &category.Category{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// MockCategories 批量创建测试分类
func MockCategories(count int) []category.Category {
	categories := make([]category.Category, count)
	for i := 0; i < count; i++ {
		categories[i] = *MockCategoryFactory(
			string(rune(i+1)),
			"分类"+string(rune(i+1)),
			"测试分类描述"+string(rune(i+1)),
		)
	}
	return categories
}

// ============== Topic Factories ==============

// MockTopicFactory 创建测试话题
func MockTopicFactory(id, title, body, userID, categoryID string) *topic.Topic {
	now := MockTime()
	return &topic.Topic{
		ID:         id,
		Title:      title,
		Body:       body,
		UserID:     userID,
		CategoryID: categoryID,
		ViewCount:  0,
		LikeCount:  0,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// MockTopics 批量创建测试话题
func MockTopics(count int) []topic.Topic {
	topics := make([]topic.Topic, count)
	for i := 0; i < count; i++ {
		topics[i] = *MockTopicFactory(
			string(rune(i+1)),
			"测试话题"+string(rune(i+1)),
			"这是测试话题的内容"+string(rune(i+1)),
			"1",
			"1",
		)
	}
	return topics
}

// ============== Comment Factories ==============

// MockCommentFactory 创建测试评论
func MockCommentFactory(id, content, userID, topicID string) *comment.Comment {
	now := MockTime()
	return &comment.Comment{
		ID:        id,
		Content:   content,
		UserID:    userID,
		TopicID:   topicID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// MockComments 批量创建测试评论
func MockComments(count int) []comment.Comment {
	comments := make([]comment.Comment, count)
	for i := 0; i < count; i++ {
		comments[i] = *MockCommentFactory(
			string(rune(i+1)),
			"测试评论内容"+string(rune(i+1)),
			"1",
			"1",
		)
	}
	return comments
}

// ============== Link Factories ==============

// MockLinkFactory 创建测试链接
func MockLinkFactory(id, title, url string) *link.Link {
	now := MockTime()
	return &link.Link{
		ID:        id,
		Title:     title,
		URL:       url,
		Sort:      0,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// MockLinks 批量创建测试链接
func MockLinks(count int) []link.Link {
	links := make([]link.Link, count)
	for i := 0; i < count; i++ {
		links[i] = *MockLinkFactory(
			string(rune(i+1)),
			"测试链接"+string(rune(i+1)),
			"https://example.com/"+string(rune(i+1)),
		)
	}
	return links
}

// ============== Role Factories ==============

// MockRoleFactory 创建测试角色
func MockRoleFactory(id, name, displayName string) *role.Role {
	now := MockTime()
	return &role.Role{
		ID:          id,
		Name:        name,
		DisplayName: displayName,
		Description: "测试角色描述",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// MockRoles 批量创建测试角色
func MockRoles(count int) []role.Role {
	roles := make([]role.Role, count)
	for i := 0; i < count; i++ {
		roles[i] = *MockRoleFactory(
			string(rune(i+1)),
			"role_"+string(rune(i+1)),
			"角色"+string(rune(i+1)),
		)
	}
	return roles
}

// ============== Permission Factories ==============

// MockPermissionFactory 创建测试权限
func MockPermissionFactory(id, name, displayName string) *permission.Permission {
	now := MockTime()
	return &permission.Permission{
		ID:          id,
		Name:        name,
		DisplayName: displayName,
		Description: "测试权限描述",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// MockPermissions 批量创建测试权限
func MockPermissions(count int) []permission.Permission {
	perms := make([]permission.Permission, count)
	for i := 0; i < count; i++ {
		perms[i] = *MockPermissionFactory(
			string(rune(i+1)),
			"permission_"+string(rune(i+1)),
			"权限"+string(rune(i+1)),
		)
	}
	return perms
}

// ============== Message Factories ==============

// MockMessageFactory 创建测试消息
func MockMessageFactory(id, content, fromUserID, toUserID string) *message.Message {
	now := MockTime()
	isRead := false
	return &message.Message{
		ID:         id,
		Content:    content,
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		IsRead:     &isRead,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// MockMessages 批量创建测试消息
func MockMessages(count int) []message.Message {
	messages := make([]message.Message, count)
	for i := 0; i < count; i++ {
		messages[i] = *MockMessageFactory(
			string(rune(i+1)),
			"测试消息内容"+string(rune(i+1)),
			"1",
			"2",
		)
	}
	return messages
}

// ============== Notification Factories ==============

// MockNotificationFactory 创建测试通知
func MockNotificationFactory(id, content, userID string) *notification.Notification {
	now := MockTime()
	isRead := false
	return &notification.Notification{
		ID:        id,
		Content:   content,
		UserID:    userID,
		Type:      "system",
		IsRead:    &isRead,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// MockNotifications 批量创建测试通知
func MockNotifications(count int) []notification.Notification {
	notifications := make([]notification.Notification, count)
	for i := 0; i < count; i++ {
		notifications[i] = *MockNotificationFactory(
			string(rune(i+1)),
			"测试通知内容"+string(rune(i+1)),
			"1",
		)
	}
	return notifications
}

// ============== Time Factories ==============

// MockTimeAfter 返回指定时间之后的时间
func MockTimeAfter(duration time.Duration) time.Time {
	return MockTime().Add(duration)
}

// MockTimeBefore 返回指定时间之前的时间
func MockTimeBefore(duration time.Duration) time.Time {
	return MockTime().Add(-duration)
}
