package services

import (
	"fmt"
	"testing"

	"GoHub-Service/app/models"
	"GoHub-Service/app/models/topic"
	"GoHub-Service/pkg/paginator"

	"github.com/gin-gonic/gin"
)

// MockTopicRepository 模拟 TopicRepository
type MockTopicRepository struct {
	topics map[string]*topic.Topic
}

func NewMockTopicRepository() *MockTopicRepository {
	return &MockTopicRepository{
		topics: make(map[string]*topic.Topic),
	}
}

func (m *MockTopicRepository) GetByID(id string) (*topic.Topic, error) {
	if t, ok := m.topics[id]; ok {
		return t, nil
	}
	return nil, nil
}

func (m *MockTopicRepository) List(c *gin.Context, perPage int) ([]topic.Topic, *paginator.Paging, error) {
	return []topic.Topic{}, &paginator.Paging{}, nil
}

func (m *MockTopicRepository) Create(t *topic.Topic) error {
	if t == nil {
		return fmt.Errorf("topic is nil")
	}
	id := fmt.Sprintf("%d", len(m.topics)+1)
	t.ID = uint64(len(m.topics) + 1)
	m.topics[id] = t
	return nil
}

func (m *MockTopicRepository) Update(t *topic.Topic) error {
	return nil
}

func (m *MockTopicRepository) Delete(id string) error {
	delete(m.topics, id)
	return nil
}

func (m *MockTopicRepository) BatchCreate(topics []topic.Topic) error {
	for i := range topics {
		_ = m.Create(&topics[i])
	}
	return nil
}

func (m *MockTopicRepository) BatchDelete(ids []string) error {
	for _, id := range ids {
		delete(m.topics, id)
	}
	return nil
}

func (m *MockTopicRepository) GetByUserID(userID string) ([]topic.Topic, error) {
	return nil, nil
}

func (m *MockTopicRepository) GetFromCache(id string) (*topic.Topic, error) {
	return nil, nil
}

func (m *MockTopicRepository) SetCache(t *topic.Topic) error {
	return nil
}

func (m *MockTopicRepository) DeleteCache(id string) error {
	return nil
}

func (m *MockTopicRepository) FlushListCache() error {
	return nil
}

// TestTopicService_Create 测试创建话题
func TestTopicService_Create(t *testing.T) {
	// 使用 Mock Repository
	mockRepo := NewMockTopicRepository()
	service := &TopicService{
		repo: mockRepo,
	}

	// 测试数据
	dto := TopicCreateDTO{
		Title:      "Test Topic",
		Body:       "Test Content",
		CategoryID: "1",
		UserID:     "1",
	}

	// 执行创建
	result, err := service.Create(dto)

	// 断言
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Error("Expected result, got nil")
	}

	if result.Title != dto.Title {
		t.Errorf("Expected title %s, got %s", dto.Title, result.Title)
	}

	if result.ID == "" {
		t.Error("Expected ID to be set")
	}
}

// TestTopicService_CheckOwnership 测试所有权检查
func TestTopicService_CheckOwnership(t *testing.T) {
	mockRepo := NewMockTopicRepository()
	service := &TopicService{
		repo: mockRepo,
	}

	// 准备测试数据
	testTopic := &topic.Topic{
		BaseModel: models.BaseModel{ID: 1},
		UserID:    "user1",
		Title:     "Test",
	}
	mockRepo.topics["1"] = testTopic

	// 测试所有者
	isOwner, err := service.CheckOwnership("1", "user1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !isOwner {
		t.Error("Expected true for owner")
	}

	// 测试非所有者
	isOwner, err = service.CheckOwnership("1", "user2")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if isOwner {
		t.Error("Expected false for non-owner")
	}
}
