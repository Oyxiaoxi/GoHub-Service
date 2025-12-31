// Package repositories User仓储测试
package repositories

import (
	"GoHub-Service/app/models/user"
	"GoHub-Service/pkg/testutil"
	"errors"
	"testing"
)

// MockUserRepositoryForTest 用于测试的用户仓储Mock
type MockUserRepositoryForTest struct {
	GetByIDFunc    func(id string) (*user.User, error)
	GetByEmailFunc func(email string) (*user.User, error)
	GetByPhoneFunc func(phone string) (*user.User, error)
	CreateFunc     func(u *user.User) error
	UpdateFunc     func(u *user.User) error
	DeleteFunc     func(id string) error
}

func (m *MockUserRepositoryForTest) GetByID(id string) (*user.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	if id == "" || id == "999" {
		return nil, errors.New("用户不存在")
	}
	return testutil.MockUserFactory(id, "测试用户", "test@example.com"), nil
}

func (m *MockUserRepositoryForTest) GetByEmail(email string) (*user.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(email)
	}
	if email == "" {
		return nil, errors.New("邮箱不能为空")
	}
	return testutil.MockUserFactory("1", "测试用户", email), nil
}

func (m *MockUserRepositoryForTest) GetByPhone(phone string) (*user.User, error) {
	if m.GetByPhoneFunc != nil {
		return m.GetByPhoneFunc(phone)
	}
	if phone == "" {
		return nil, errors.New("手机号不能为空")
	}
	u := testutil.MockUserFactory("1", "测试用户", "test@example.com")
	u.Phone = phone
	return u, nil
}

func (m *MockUserRepositoryForTest) Create(u *user.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(u)
	}
	if u == nil {
		return errors.New("用户不能为nil")
	}
	if u.Email == "" {
		return errors.New("邮箱不能为空")
	}
	return nil
}

func (m *MockUserRepositoryForTest) Update(u *user.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(u)
	}
	if u == nil {
		return errors.New("用户不能为nil")
	}
	return nil
}

func (m *MockUserRepositoryForTest) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	if id == "" || id == "999" {
		return errors.New("用户不存在")
	}
	return nil
}

func (m *MockUserRepositoryForTest) List(c interface{}, perPage int) ([]user.User, interface{}, error) {
	return []user.User{
		*testutil.MockUserFactory("1", "用户1", "user1@example.com"),
		*testutil.MockUserFactory("2", "用户2", "user2@example.com"),
	}, nil, nil
}

func (m *MockUserRepositoryForTest) IncrementNotificationCount(userID string) error {
	if userID == "" || userID == "999" {
		return errors.New("用户不存在")
	}
	return nil
}

func (m *MockUserRepositoryForTest) ClearNotificationCount(userID string) error {
	if userID == "" || userID == "999" {
		return errors.New("用户不存在")
	}
	return nil
}

func (m *MockUserRepositoryForTest) UpdateLastActiveAt(userID string) error {
	if userID == "" || userID == "999" {
		return errors.New("用户不存在")
	}
	return nil
}

// TestMockUserRepository_GetByID 测试通过ID获取用户
func TestMockUserRepository_GetByID(t *testing.T) {
	mockRepo := &MockUserRepositoryForTest{}

	tests := []struct {
		name      string
		id        string
		wantErr   bool
		checkFunc func(t *testing.T, result *user.User)
	}{
		{
			name:    "成功获取用户",
			id:      "1",
			wantErr: false,
			checkFunc: func(t *testing.T, result *user.User) {
				testutil.AssertNotNil(t, result, "用户不应为nil")
				testutil.AssertEqual(t, "1", result.ID, "用户ID应该匹配")
			},
		},
		{
			name:    "用户不存在",
			id:      "999",
			wantErr: true,
		},
		{
			name:    "空ID",
			id:      "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := mockRepo.GetByID(tt.id)

			if tt.wantErr {
				testutil.AssertNotNil(t, err, "应该返回错误")
			} else {
				testutil.AssertNil(t, err, "不应该返回错误")
				if tt.checkFunc != nil {
					tt.checkFunc(t, result)
				}
			}
		})
	}
}

// TestMockUserRepository_GetByEmail 测试通过邮箱获取用户
func TestMockUserRepository_GetByEmail(t *testing.T) {
	mockRepo := &MockUserRepositoryForTest{}

	tests := []struct {
		name      string
		email     string
		wantErr   bool
		checkFunc func(t *testing.T, result *user.User)
	}{
		{
			name:    "成功获取用户",
			email:   "test@example.com",
			wantErr: false,
			checkFunc: func(t *testing.T, result *user.User) {
				testutil.AssertNotNil(t, result, "用户不应为nil")
				testutil.AssertEqual(t, "test@example.com", result.Email, "邮箱应该匹配")
			},
		},
		{
			name:    "空邮箱",
			email:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := mockRepo.GetByEmail(tt.email)

			if tt.wantErr {
				testutil.AssertNotNil(t, err, "应该返回错误")
			} else {
				testutil.AssertNil(t, err, "不应该返回错误")
				if tt.checkFunc != nil {
					tt.checkFunc(t, result)
				}
			}
		})
	}
}

// TestMockUserRepository_Create 测试创建用户
func TestMockUserRepository_Create(t *testing.T) {
	mockRepo := &MockUserRepositoryForTest{}

	tests := []struct {
		name    string
		user    *user.User
		wantErr bool
	}{
		{
			name:    "成功创建用户",
			user:    testutil.MockUserFactory("1", "新用户", "newuser@example.com"),
			wantErr: false,
		},
		{
			name:    "nil用户",
			user:    nil,
			wantErr: true,
		},
		{
			name: "邮箱为空",
			user: &user.User{
				Name:  "用户",
				Email: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mockRepo.Create(tt.user)

			if tt.wantErr {
				testutil.AssertNotNil(t, err, "应该返回错误")
			} else {
				testutil.AssertNil(t, err, "不应该返回错误")
			}
		})
	}
}

// TestMockUserRepository_Update 测试更新用户
func TestMockUserRepository_Update(t *testing.T) {
	mockRepo := &MockUserRepositoryForTest{}

	tests := []struct {
		name    string
		user    *user.User
		wantErr bool
	}{
		{
			name:    "成功更新用户",
			user:    testutil.MockUserFactory("1", "更新用户", "update@example.com"),
			wantErr: false,
		},
		{
			name:    "nil用户",
			user:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mockRepo.Update(tt.user)

			if tt.wantErr {
				testutil.AssertNotNil(t, err, "应该返回错误")
			} else {
				testutil.AssertNil(t, err, "不应该返回错误")
			}
		})
	}
}

// TestMockUserRepository_Delete 测试删除用户
func TestMockUserRepository_Delete(t *testing.T) {
	mockRepo := &MockUserRepositoryForTest{}

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "成功删除用户",
			id:      "1",
			wantErr: false,
		},
		{
			name:    "用户不存在",
			id:      "999",
			wantErr: true,
		},
		{
			name:    "空ID",
			id:      "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mockRepo.Delete(tt.id)

			if tt.wantErr {
				testutil.AssertNotNil(t, err, "应该返回错误")
			} else {
				testutil.AssertNil(t, err, "不应该返回错误")
			}
		})
	}
}

// TestMockUserRepository_IncrementNotificationCount 测试增加通知计数
func TestMockUserRepository_IncrementNotificationCount(t *testing.T) {
	mockRepo := &MockUserRepositoryForTest{}

	tests := []struct {
		name    string
		userID  string
		wantErr bool
	}{
		{
			name:    "成功增加通知计数",
			userID:  "1",
			wantErr: false,
		},
		{
			name:    "用户不存在",
			userID:  "999",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mockRepo.IncrementNotificationCount(tt.userID)

			if tt.wantErr {
				testutil.AssertNotNil(t, err, "应该返回错误")
			} else {
				testutil.AssertNil(t, err, "不应该返回错误")
			}
		})
	}
}

// TestMockUserRepository_List 测试获取用户列表
func TestMockUserRepository_List(t *testing.T) {
	mockRepo := &MockUserRepositoryForTest{}

	users, paging, err := mockRepo.List(nil, 10)

	testutil.AssertNil(t, err, "不应该返回错误")
	testutil.AssertNotNil(t, users, "用户列表不应为nil")
	testutil.AssertEqual(t, 2, len(users), "应该返回2个用户")
	testutil.AssertNil(t, paging, "分页信息应为nil（Mock实现）")
}
