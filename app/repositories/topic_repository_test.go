package repositories

import (
	"GoHub-Service/app/models/topic"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockTopicRepository 用于测试的 Mock 实现
type MockTopicRepository struct {
	GetByIDFunc    func(id string) (*topic.Topic, error)
	ListFunc       func(perPage int) ([]topic.Topic, error)
	CreateFunc     func(t *topic.Topic) error
	UpdateFunc     func(t *topic.Topic) error
	DeleteFunc     func(id string) error
	BatchCreateFunc func(topics []topic.Topic) error
	BatchDeleteFunc func(ids []string) error
}

func (m *MockTopicRepository) GetByID(id string) (*topic.Topic, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *MockTopicRepository) List(c interface{}, perPage int) ([]topic.Topic, interface{}, error) {
	if m.ListFunc != nil {
		topics, err := m.ListFunc(perPage)
		return topics, nil, err
	}
	return []topic.Topic{}, nil, nil
}

func (m *MockTopicRepository) Create(t *topic.Topic) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(t)
	}
	return nil
}

func (m *MockTopicRepository) Update(t *topic.Topic) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(t)
	}
	return nil
}

func (m *MockTopicRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func (m *MockTopicRepository) BatchCreate(topics []topic.Topic) error {
	if m.BatchCreateFunc != nil {
		return m.BatchCreateFunc(topics)
	}
	return nil
}

func (m *MockTopicRepository) BatchDelete(ids []string) error {
	if m.BatchDeleteFunc != nil {
		return m.BatchDeleteFunc(ids)
	}
	return nil
}

// TestMockTopicRepository_GetByID 测试 GetByID Mock
func TestMockTopicRepository_GetByID(t *testing.T) {
	mock := &MockTopicRepository{
		GetByIDFunc: func(id string) (*topic.Topic, error) {
			if id == "1" {
				return &topic.Topic{Title: "Test Topic"}, nil
			}
			return nil, ErrNotFound
		},
	}

	// 测试存在的话题
	result, err := mock.GetByID("1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Topic", result.Title)

	// 测试不存在的话题
	result, err = mock.GetByID("999")
	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestMockTopicRepository_Create 测试 Create Mock
func TestMockTopicRepository_Create(t *testing.T) {
	mock := &MockTopicRepository{
		CreateFunc: func(tp *topic.Topic) error {
			if tp.Title == "" {
				return ErrCreateFailed
			}
			tp.ID = 1
			return nil
		},
	}

	// 测试成功创建
	newTopic := &topic.Topic{Title: "New Topic"}
	err := mock.Create(newTopic)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), newTopic.ID)

	// 测试创建失败
	invalidTopic := &topic.Topic{Title: ""}
	err = mock.Create(invalidTopic)
	assert.Error(t, err)
}

// TestMockTopicRepository_Update 测试 Update Mock
func TestMockTopicRepository_Update(t *testing.T) {
	mock := &MockTopicRepository{
		UpdateFunc: func(tp *topic.Topic) error {
			if tp.ID == 0 {
				return ErrUpdateFailed
			}
			return nil
		},
	}

	// 测试成功更新
	existingTopic := &topic.Topic{Title: "Updated Topic"}
	existingTopic.ID = 1
	err := mock.Update(existingTopic)
	assert.NoError(t, err)

	// 测试更新失败（ID 为 0）
	invalidTopic := &topic.Topic{Title: "Invalid"}
	err = mock.Update(invalidTopic)
	assert.Error(t, err)
}

// TestMockTopicRepository_Delete 测试 Delete Mock
func TestMockTopicRepository_Delete(t *testing.T) {
	mock := &MockTopicRepository{
		DeleteFunc: func(id string) error {
			if id == "1" {
				return nil
			}
			return ErrNotFound
		},
	}

	// 测试成功删除
	err := mock.Delete("1")
	assert.NoError(t, err)

	// 测试删除不存在的话题
	err = mock.Delete("999")
	assert.Error(t, err)
}

// TestMockTopicRepository_BatchCreate 测试 BatchCreate Mock
func TestMockTopicRepository_BatchCreate(t *testing.T) {
	mock := &MockTopicRepository{
		BatchCreateFunc: func(topics []topic.Topic) error {
			if len(topics) == 0 {
				return ErrCreateFailed
			}
			for i := range topics {
				topics[i].ID = uint64(i + 1)
			}
			return nil
		},
	}

	// 测试批量创建成功
	topics := []topic.Topic{
		{Title: "Topic 1"},
		{Title: "Topic 2"},
	}
	err := mock.BatchCreate(topics)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), topics[0].ID)
	assert.Equal(t, uint64(2), topics[1].ID)

	// 测试批量创建失败（空数组）
	err = mock.BatchCreate([]topic.Topic{})
	assert.Error(t, err)
}

// TestMockTopicRepository_BatchDelete 测试 BatchDelete Mock
func TestMockTopicRepository_BatchDelete(t *testing.T) {
	mock := &MockTopicRepository{
		BatchDeleteFunc: func(ids []string) error {
			if len(ids) == 0 {
				return ErrDeleteFailed
			}
			return nil
		},
	}

	// 测试批量删除成功
	err := mock.BatchDelete([]string{"1", "2"})
	assert.NoError(t, err)

	// 测试批量删除失败（空数组）
	err = mock.BatchDelete([]string{})
	assert.Error(t, err)
}
