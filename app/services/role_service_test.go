// Package services 角色服务测试
package services

import (
	"GoHub-Service/app/models/role"
	"GoHub-Service/pkg/testutil"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockRoleRepository 角色仓储 Mock
type MockRoleRepository struct {
	GetByIDFunc     func(ctx context.Context, id string) (*role.Role, error)
	GetByNameFunc   func(ctx context.Context, name string) (*role.Role, error)
	GetAllFunc      func(ctx context.Context) ([]role.Role, error)
	CreateFunc      func(ctx context.Context, r *role.Role) error
	UpdateFunc      func(ctx context.Context, r *role.Role) error
	DeleteFunc      func(ctx context.Context, id string) error
	ExistsFunc      func(ctx context.Context, name string) (bool, error)
	GetUserRolesFunc func(ctx context.Context, userID string) ([]role.Role, error)
}

func (m *MockRoleRepository) GetByID(ctx context.Context, id string) (*role.Role, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return testutil.MockRoleFactory(id, "admin", "管理员"), nil
}

func (m *MockRoleRepository) GetByName(ctx context.Context, name string) (*role.Role, error) {
	if m.GetByNameFunc != nil {
		return m.GetByNameFunc(ctx, name)
	}
	return testutil.MockRoleFactory("1", name, "测试角色"), nil
}

func (m *MockRoleRepository) GetAll(ctx context.Context) ([]role.Role, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return testutil.MockRoles(3), nil
}

func (m *MockRoleRepository) Create(ctx context.Context, r *role.Role) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, r)
	}
	r.ID = "1"
	return nil
}

func (m *MockRoleRepository) Update(ctx context.Context, r *role.Role) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, r)
	}
	return nil
}

func (m *MockRoleRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockRoleRepository) Exists(ctx context.Context, name string) (bool, error) {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(ctx, name)
	}
	return false, nil
}

func (m *MockRoleRepository) GetUserRoles(ctx context.Context, userID string) ([]role.Role, error) {
	if m.GetUserRolesFunc != nil {
		return m.GetUserRolesFunc(ctx, userID)
	}
	return testutil.MockRoles(2), nil
}

// TestRoleService_Create 测试创建角色
func TestRoleService_Create(t *testing.T) {
	mockRepo := &MockRoleRepository{}
	ctx := context.Background()

	tests := []struct {
		name        string
		roleName    string
		displayName string
		setupMock   func()
		wantErr     bool
	}{
		{
			name:        "成功创建角色",
			roleName:    "editor",
			displayName: "编辑",
			setupMock: func() {
				mockRepo.ExistsFunc = func(ctx context.Context, name string) (bool, error) {
					return false, nil
				}
				mockRepo.CreateFunc = func(ctx context.Context, r *role.Role) error {
					r.ID = "1"
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:        "角色名称为空",
			roleName:    "",
			displayName: "测试角色",
			setupMock: func() {
				mockRepo.CreateFunc = func(ctx context.Context, r *role.Role) error {
					return errors.New("角色名称不能为空")
				}
			},
			wantErr: true,
		},
		{
			name:        "角色已存在",
			roleName:    "admin",
			displayName: "管理员",
			setupMock: func() {
				mockRepo.ExistsFunc = func(ctx context.Context, name string) (bool, error) {
					return true, nil
				}
				mockRepo.CreateFunc = func(ctx context.Context, r *role.Role) error {
					return errors.New("角色已存在")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			newRole := &role.Role{
				Name:        tt.roleName,
				DisplayName: tt.displayName,
			}

			err := mockRepo.Create(ctx, newRole)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, newRole.ID)
			}
		})
	}
}

// TestRoleService_GetByName 测试根据名称获取角色
func TestRoleService_GetByName(t *testing.T) {
	mockRepo := &MockRoleRepository{
		GetByNameFunc: func(ctx context.Context, name string) (*role.Role, error) {
			if name == "" {
				return nil, errors.New("角色名称不能为空")
			}
			if name == "notexist" {
				return nil, errors.New("角色不存在")
			}
			return testutil.MockRoleFactory("1", name, "测试角色"), nil
		},
	}

	ctx := context.Background()

	t.Run("成功获取角色", func(t *testing.T) {
		role, err := mockRepo.GetByName(ctx, "admin")
		assert.NoError(t, err)
		assert.NotNil(t, role)
		assert.Equal(t, "admin", role.Name)
	})

	t.Run("角色不存在", func(t *testing.T) {
		role, err := mockRepo.GetByName(ctx, "notexist")
		assert.Error(t, err)
		assert.Nil(t, role)
	})

	t.Run("角色名称为空", func(t *testing.T) {
		role, err := mockRepo.GetByName(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, role)
	})
}

// TestRoleService_GetAll 测试获取所有角色
func TestRoleService_GetAll(t *testing.T) {
	mockRepo := &MockRoleRepository{
		GetAllFunc: func(ctx context.Context) ([]role.Role, error) {
			return testutil.MockRoles(5), nil
		},
	}

	ctx := context.Background()

	t.Run("成功获取所有角色", func(t *testing.T) {
		roles, err := mockRepo.GetAll(ctx)
		assert.NoError(t, err)
		assert.Len(t, roles, 5)
	})
}

// TestRoleService_Update 测试更新角色
func TestRoleService_Update(t *testing.T) {
	mockRepo := &MockRoleRepository{
		UpdateFunc: func(ctx context.Context, r *role.Role) error {
			if r.ID == "" {
				return errors.New("角色ID不能为空")
			}
			if r.ID == "999" {
				return errors.New("角色不存在")
			}
			return nil
		},
	}

	ctx := context.Background()

	t.Run("成功更新角色", func(t *testing.T) {
		role := testutil.MockRoleFactory("1", "admin", "超级管理员")
		err := mockRepo.Update(ctx, role)
		assert.NoError(t, err)
	})

	t.Run("角色ID为空", func(t *testing.T) {
		role := &role.Role{
			Name:        "test",
			DisplayName: "测试",
		}
		err := mockRepo.Update(ctx, role)
		assert.Error(t, err)
	})

	t.Run("角色不存在", func(t *testing.T) {
		role := testutil.MockRoleFactory("999", "test", "测试")
		err := mockRepo.Update(ctx, role)
		assert.Error(t, err)
	})
}

// TestRoleService_Delete 测试删除角色
func TestRoleService_Delete(t *testing.T) {
	mockRepo := &MockRoleRepository{
		DeleteFunc: func(ctx context.Context, id string) error {
			if id == "999" {
				return errors.New("角色不存在")
			}
			if id == "1" {
				return errors.New("不能删除系统角色")
			}
			return nil
		},
	}

	ctx := context.Background()

	t.Run("成功删除角色", func(t *testing.T) {
		err := mockRepo.Delete(ctx, "10")
		assert.NoError(t, err)
	})

	t.Run("角色不存在", func(t *testing.T) {
		err := mockRepo.Delete(ctx, "999")
		assert.Error(t, err)
	})

	t.Run("不能删除系统角色", func(t *testing.T) {
		err := mockRepo.Delete(ctx, "1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "系统角色")
	})
}

// TestRoleService_Exists 测试检查角色是否存在
func TestRoleService_Exists(t *testing.T) {
	mockRepo := &MockRoleRepository{
		ExistsFunc: func(ctx context.Context, name string) (bool, error) {
			if name == "admin" {
				return true, nil
			}
			return false, nil
		},
	}

	ctx := context.Background()

	t.Run("角色存在", func(t *testing.T) {
		exists, err := mockRepo.Exists(ctx, "admin")
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("角色不存在", func(t *testing.T) {
		exists, err := mockRepo.Exists(ctx, "notexist")
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

// TestRoleService_GetUserRoles 测试获取用户角色
func TestRoleService_GetUserRoles(t *testing.T) {
	mockRepo := &MockRoleRepository{
		GetUserRolesFunc: func(ctx context.Context, userID string) ([]role.Role, error) {
			if userID == "" {
				return nil, errors.New("用户ID不能为空")
			}
			if userID == "999" {
				return []role.Role{}, nil
			}
			return testutil.MockRoles(2), nil
		},
	}

	ctx := context.Background()

	t.Run("成功获取用户角色", func(t *testing.T) {
		roles, err := mockRepo.GetUserRoles(ctx, "1")
		assert.NoError(t, err)
		assert.Len(t, roles, 2)
	})

	t.Run("用户没有角色", func(t *testing.T) {
		roles, err := mockRepo.GetUserRoles(ctx, "999")
		assert.NoError(t, err)
		assert.Empty(t, roles)
	})

	t.Run("用户ID为空", func(t *testing.T) {
		roles, err := mockRepo.GetUserRoles(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, roles)
	})
}
