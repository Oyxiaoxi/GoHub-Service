// Package services 用户服务测试
package services

import (
	"GoHub-Service/app/models/user"
	"GoHub-Service/app/repositories"
	"GoHub-Service/pkg/testutil"
	"errors"
	"testing"
)

// MockUserRepository 用户仓储Mock
type MockUserRepository struct {
	GetByIDFunc     func(id string) (*user.User, error)
	GetByEmailFunc  func(email string) (*user.User, error)
	GetByPhoneFunc  func(phone string) (*user.User, error)
	CreateFunc      func(u *user.User) error
	UpdateFunc      func(u *user.User) error
	DeleteFunc      func(id string) error
	BatchCreateFunc func(users []user.User) error
	BatchDeleteFunc func(ids []string) error
}

func (m *MockUserRepository) GetByID(id string) (*user.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return testutil.MockUserFactory("1", "测试用户", "test@example.com"), nil
}

func (m *MockUserRepository) GetByEmail(email string) (*user.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(email)
	}
	return nil, errors.New("用户不存在")
}

func (m *MockUserRepository) GetByPhone(phone string) (*user.User, error) {
	if m.GetByPhoneFunc != nil {
		return m.GetByPhoneFunc(phone)
	}
	return nil, errors.New("用户不存在")
}

func (m *MockUserRepository) Create(u *user.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(u)
	}
	return nil
}

func (m *MockUserRepository) Update(u *user.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(u)
	}
	return nil
}

func (m *MockUserRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func (m *MockUserRepository) BatchCreate(users []user.User) error {
	if m.BatchCreateFunc != nil {
		return m.BatchCreateFunc(users)
	}
	return nil
}

func (m *MockUserRepository) BatchDelete(ids []string) error {
	if m.BatchDeleteFunc != nil {
		return m.BatchDeleteFunc(ids)
	}
	return nil
}

func (m *MockUserRepository) List(c interface{}, perPage int) ([]user.User, interface{}, error) {
	return []user.User{}, nil, nil
}

func (m *MockUserRepository) IncrementNotificationCount(userID string) error {
	return nil
}

func (m *MockUserRepository) ClearNotificationCount(userID string) error {
	return nil
}

func (m *MockUserRepository) UpdateLastActiveAt(userID string) error {
	return nil
}

func (m *MockUserRepository) GetFromCache(id string) (*user.User, error) {
	return nil, nil
}

func (m *MockUserRepository) SetCache(u *user.User) error {
	return nil
}

func (m *MockUserRepository) DeleteCache(id string) error {
	return nil
}

// 确保MockUserRepository实现了UserRepository接口
var _ repositories.UserRepository = (*MockUserRepository)(nil)

// TestUserService_GetByID 测试通过ID获取用户
func TestUserService_GetByID(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		mockFunc  func(id string) (*user.User, error)
		wantErr   bool
		checkFunc func(t *testing.T, result *user.User)
	}{
		{
			name: "成功获取用户",
			id:   "1",
			mockFunc: func(id string) (*user.User, error) {
				return testutil.MockUserFactory("1", "张三", "zhangsan@example.com"), nil
			},
			wantErr: false,
			checkFunc: func(t *testing.T, result *user.User) {
				testutil.AssertNotNil(t, result, "用户不应为nil")
				testutil.AssertEqual(t, "1", result.ID, "用户ID应该匹配")
				testutil.AssertEqual(t, "张三", result.Name, "用户名应该匹配")
				testutil.AssertEqual(t, "zhangsan@example.com", result.Email, "邮箱应该匹配")
			},
		},
		{
			name: "用户不存在",
			id:   "999",
			mockFunc: func(id string) (*user.User, error) {
				return nil, errors.New("用户不存在")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{
				GetByIDFunc: tt.mockFunc,
			}

			service := &UserService{repo: mockRepo}
			result, err := service.GetByID(tt.id)

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

// TestUserService_GetByEmail 测试通过邮箱获取用户
func TestUserService_GetByEmail(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		mockFunc  func(email string) (*user.User, error)
		wantErr   bool
		checkFunc func(t *testing.T, result *user.User)
	}{
		{
			name:  "成功获取用户",
			email: "test@example.com",
			mockFunc: func(email string) (*user.User, error) {
				return testutil.MockUserFactory("1", "测试用户", email), nil
			},
			wantErr: false,
			checkFunc: func(t *testing.T, result *user.User) {
				testutil.AssertNotNil(t, result, "用户不应为nil")
				testutil.AssertEqual(t, "test@example.com", result.Email, "邮箱应该匹配")
			},
		},
		{
			name:  "邮箱不存在",
			email: "notexist@example.com",
			mockFunc: func(email string) (*user.User, error) {
				return nil, errors.New("用户不存在")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{
				GetByEmailFunc: tt.mockFunc,
			}

			service := &UserService{repo: mockRepo}
			result, err := service.GetByEmail(tt.email)

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

// TestUserService_GetByPhone 测试通过手机号获取用户
func TestUserService_GetByPhone(t *testing.T) {
	tests := []struct {
		name      string
		phone     string
		mockFunc  func(phone string) (*user.User, error)
		wantErr   bool
		checkFunc func(t *testing.T, result *user.User)
	}{
		{
			name:  "成功获取用户",
			phone: "13800138000",
			mockFunc: func(phone string) (*user.User, error) {
				u := testutil.MockUserFactory("1", "测试用户", "test@example.com")
				u.Phone = phone
				return u, nil
			},
			wantErr: false,
			checkFunc: func(t *testing.T, result *user.User) {
				testutil.AssertNotNil(t, result, "用户不应为nil")
				testutil.AssertEqual(t, "13800138000", result.Phone, "手机号应该匹配")
			},
		},
		{
			name:  "手机号不存在",
			phone: "18800000000",
			mockFunc: func(phone string) (*user.User, error) {
				return nil, errors.New("用户不存在")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{
				GetByPhoneFunc: tt.mockFunc,
			}

			service := &UserService{repo: mockRepo}
			result, err := service.GetByPhone(tt.phone)

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

// TestUserService_IncrementNotificationCount 测试增加通知计数
func TestUserService_IncrementNotificationCount(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		mockFunc func(userID string) error
		wantErr  bool
	}{
		{
			name:   "成功增加通知计数",
			userID: "1",
			mockFunc: func(userID string) error {
				return nil
			},
			wantErr: false,
		},
		{
			name:   "用户不存在",
			userID: "999",
			mockFunc: func(userID string) error {
				return errors.New("用户不存在")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{
				IncrementNotificationCount: func(userID string) error {
					if tt.mockFunc != nil {
						return tt.mockFunc(userID)
					}
					return nil
				},
			}

			service := &UserService{repo: mockRepo}
			err := service.IncrementNotificationCount(tt.userID)

			if tt.wantErr {
				testutil.AssertNotNil(t, err, "应该返回错误")
			} else {
				testutil.AssertNil(t, err, "不应该返回错误")
			}
		})
	}
}

// TestUserService_ClearNotificationCount 测试清空通知计数
func TestUserService_ClearNotificationCount(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		mockFunc func(userID string) error
		wantErr  bool
	}{
		{
			name:   "成功清空通知计数",
			userID: "1",
			mockFunc: func(userID string) error {
				return nil
			},
			wantErr: false,
		},
		{
			name:   "用户不存在",
			userID: "999",
			mockFunc: func(userID string) error {
				return errors.New("用户不存在")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{
				ClearNotificationCount: func(userID string) error {
					if tt.mockFunc != nil {
						return tt.mockFunc(userID)
					}
					return nil
				},
			}

			service := &UserService{repo: mockRepo}
			err := service.ClearNotificationCount(tt.userID)

			if tt.wantErr {
				testutil.AssertNotNil(t, err, "应该返回错误")
			} else {
				testutil.AssertNil(t, err, "不应该返回错误")
			}
		})
	}
}

// TestUserService_UpdateLastActiveAt 测试更新最后活跃时间
func TestUserService_UpdateLastActiveAt(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		mockFunc func(userID string) error
		wantErr  bool
	}{
		{
			name:   "成功更新最后活跃时间",
			userID: "1",
			mockFunc: func(userID string) error {
				return nil
			},
			wantErr: false,
		},
		{
			name:   "更新失败",
			userID: "999",
			mockFunc: func(userID string) error {
				return errors.New("更新失败")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{
				UpdateLastActiveAt: func(userID string) error {
					if tt.mockFunc != nil {
						return tt.mockFunc(userID)
					}
					return nil
				},
			}

			service := &UserService{repo: mockRepo}
			err := service.UpdateLastActiveAt(tt.userID)

			if tt.wantErr {
				testutil.AssertNotNil(t, err, "应该返回错误")
			} else {
				testutil.AssertNil(t, err, "不应该返回错误")
			}
		})
	}
}
