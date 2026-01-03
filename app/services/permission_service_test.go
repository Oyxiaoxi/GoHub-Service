// Package services 权限服务测试
package services

import (
	"GoHub-Service/app/models/permission"
	"GoHub-Service/pkg/testutil"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockPermissionRepository 权限仓储 Mock
type MockPermissionRepository struct {
	GetByIDFunc          func(ctx context.Context, id string) (*permission.Permission, error)
	GetByNameFunc        func(ctx context.Context, name string) (*permission.Permission, error)
	GetAllFunc           func(ctx context.Context) ([]permission.Permission, error)
	CreateFunc           func(ctx context.Context, p *permission.Permission) error
	UpdateFunc           func(ctx context.Context, p *permission.Permission) error
	DeleteFunc           func(ctx context.Context, id string) error
	ExistsFunc           func(ctx context.Context, name string) (bool, error)
	GetRolePermissionsFunc func(ctx context.Context, roleID string) ([]permission.Permission, error)
}

func (m *MockPermissionRepository) GetByID(ctx context.Context, id string) (*permission.Permission, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return testutil.MockPermissionFactory(id, "read_post", "查看文章"), nil
}

func (m *MockPermissionRepository) GetByName(ctx context.Context, name string) (*permission.Permission, error) {
	if m.GetByNameFunc != nil {
		return m.GetByNameFunc(ctx, name)
	}
	return testutil.MockPermissionFactory("1", name, "测试权限"), nil
}

func (m *MockPermissionRepository) GetAll(ctx context.Context) ([]permission.Permission, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return testutil.MockPermissions(5), nil
}

func (m *MockPermissionRepository) Create(ctx context.Context, p *permission.Permission) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, p)
	}
	p.ID = "1"
	return nil
}

func (m *MockPermissionRepository) Update(ctx context.Context, p *permission.Permission) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, p)
	}
	return nil
}

func (m *MockPermissionRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockPermissionRepository) Exists(ctx context.Context, name string) (bool, error) {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(ctx, name)
	}
	return false, nil
}

func (m *MockPermissionRepository) GetRolePermissions(ctx context.Context, roleID string) ([]permission.Permission, error) {
	if m.GetRolePermissionsFunc != nil {
		return m.GetRolePermissionsFunc(ctx, roleID)
	}
	return testutil.MockPermissions(3), nil
}

// TestPermissionService_Create 测试创建权限
func TestPermissionService_Create(t *testing.T) {
	mockRepo := &MockPermissionRepository{}
	ctx := context.Background()

	tests := []struct {
		name            string
		permissionName  string
		displayName     string
		setupMock       func()
		wantErr         bool
	}{
		{
			name:           "成功创建权限",
			permissionName: "create_post",
			displayName:    "创建文章",
			setupMock: func() {
				mockRepo.ExistsFunc = func(ctx context.Context, name string) (bool, error) {
					return false, nil
				}
				mockRepo.CreateFunc = func(ctx context.Context, p *permission.Permission) error {
					p.ID = "1"
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:           "权限名称为空",
			permissionName: "",
			displayName:    "测试权限",
			setupMock: func() {
				mockRepo.CreateFunc = func(ctx context.Context, p *permission.Permission) error {
					return errors.New("权限名称不能为空")
				}
			},
			wantErr: true,
		},
		{
			name:           "权限已存在",
			permissionName: "read_post",
			displayName:    "查看文章",
			setupMock: func() {
				mockRepo.ExistsFunc = func(ctx context.Context, name string) (bool, error) {
					return true, nil
				}
				mockRepo.CreateFunc = func(ctx context.Context, p *permission.Permission) error {
					return errors.New("权限已存在")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			newPermission := &permission.Permission{
				Name:        tt.permissionName,
				DisplayName: tt.displayName,
			}

			err := mockRepo.Create(ctx, newPermission)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, newPermission.ID)
			}
		})
	}
}

// TestPermissionService_GetByName 测试根据名称获取权限
func TestPermissionService_GetByName(t *testing.T) {
	mockRepo := &MockPermissionRepository{
		GetByNameFunc: func(ctx context.Context, name string) (*permission.Permission, error) {
			if name == "" {
				return nil, errors.New("权限名称不能为空")
			}
			if name == "notexist" {
				return nil, errors.New("权限不存在")
			}
			return testutil.MockPermissionFactory("1", name, "测试权限"), nil
		},
	}

	ctx := context.Background()

	t.Run("成功获取权限", func(t *testing.T) {
		perm, err := mockRepo.GetByName(ctx, "read_post")
		assert.NoError(t, err)
		assert.NotNil(t, perm)
		assert.Equal(t, "read_post", perm.Name)
	})

	t.Run("权限不存在", func(t *testing.T) {
		perm, err := mockRepo.GetByName(ctx, "notexist")
		assert.Error(t, err)
		assert.Nil(t, perm)
	})

	t.Run("权限名称为空", func(t *testing.T) {
		perm, err := mockRepo.GetByName(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, perm)
	})
}

// TestPermissionService_GetAll 测试获取所有权限
func TestPermissionService_GetAll(t *testing.T) {
	mockRepo := &MockPermissionRepository{
		GetAllFunc: func(ctx context.Context) ([]permission.Permission, error) {
			return testutil.MockPermissions(10), nil
		},
	}

	ctx := context.Background()

	t.Run("成功获取所有权限", func(t *testing.T) {
		perms, err := mockRepo.GetAll(ctx)
		assert.NoError(t, err)
		assert.Len(t, perms, 10)
	})
}

// TestPermissionService_Update 测试更新权限
func TestPermissionService_Update(t *testing.T) {
	mockRepo := &MockPermissionRepository{
		UpdateFunc: func(ctx context.Context, p *permission.Permission) error {
			if p.ID == "" {
				return errors.New("权限ID不能为空")
			}
			if p.ID == "999" {
				return errors.New("权限不存在")
			}
			return nil
		},
	}

	ctx := context.Background()

	t.Run("成功更新权限", func(t *testing.T) {
		perm := testutil.MockPermissionFactory("1", "edit_post", "编辑文章")
		err := mockRepo.Update(ctx, perm)
		assert.NoError(t, err)
	})

	t.Run("权限ID为空", func(t *testing.T) {
		perm := &permission.Permission{
			Name:        "test",
			DisplayName: "测试",
		}
		err := mockRepo.Update(ctx, perm)
		assert.Error(t, err)
	})

	t.Run("权限不存在", func(t *testing.T) {
		perm := testutil.MockPermissionFactory("999", "test", "测试")
		err := mockRepo.Update(ctx, perm)
		assert.Error(t, err)
	})
}

// TestPermissionService_Delete 测试删除权限
func TestPermissionService_Delete(t *testing.T) {
	mockRepo := &MockPermissionRepository{
		DeleteFunc: func(ctx context.Context, id string) error {
			if id == "999" {
				return errors.New("权限不存在")
			}
			if id == "1" {
				return errors.New("不能删除系统权限")
			}
			return nil
		},
	}

	ctx := context.Background()

	t.Run("成功删除权限", func(t *testing.T) {
		err := mockRepo.Delete(ctx, "10")
		assert.NoError(t, err)
	})

	t.Run("权限不存在", func(t *testing.T) {
		err := mockRepo.Delete(ctx, "999")
		assert.Error(t, err)
	})

	t.Run("不能删除系统权限", func(t *testing.T) {
		err := mockRepo.Delete(ctx, "1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "系统权限")
	})
}

// TestPermissionService_Exists 测试检查权限是否存在
func TestPermissionService_Exists(t *testing.T) {
	mockRepo := &MockPermissionRepository{
		ExistsFunc: func(ctx context.Context, name string) (bool, error) {
			if name == "read_post" {
				return true, nil
			}
			return false, nil
		},
	}

	ctx := context.Background()

	t.Run("权限存在", func(t *testing.T) {
		exists, err := mockRepo.Exists(ctx, "read_post")
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("权限不存在", func(t *testing.T) {
		exists, err := mockRepo.Exists(ctx, "notexist")
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

// TestPermissionService_GetRolePermissions 测试获取角色权限
func TestPermissionService_GetRolePermissions(t *testing.T) {
	mockRepo := &MockPermissionRepository{
		GetRolePermissionsFunc: func(ctx context.Context, roleID string) ([]permission.Permission, error) {
			if roleID == "" {
				return nil, errors.New("角色ID不能为空")
			}
			if roleID == "999" {
				return []permission.Permission{}, nil
			}
			return testutil.MockPermissions(5), nil
		},
	}

	ctx := context.Background()

	t.Run("成功获取角色权限", func(t *testing.T) {
		perms, err := mockRepo.GetRolePermissions(ctx, "1")
		assert.NoError(t, err)
		assert.Len(t, perms, 5)
	})

	t.Run("角色没有权限", func(t *testing.T) {
		perms, err := mockRepo.GetRolePermissions(ctx, "999")
		assert.NoError(t, err)
		assert.Empty(t, perms)
	})

	t.Run("角色ID为空", func(t *testing.T) {
		perms, err := mockRepo.GetRolePermissions(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, perms)
	})
}

// TestPermissionService_BatchCreate 测试批量创建权限
func TestPermissionService_BatchCreate(t *testing.T) {
	mockRepo := &MockPermissionRepository{}
	ctx := context.Background()

	t.Run("批量创建权限", func(t *testing.T) {
		perms := testutil.MockPermissions(5)
		
		for i := range perms {
			mockRepo.CreateFunc = func(ctx context.Context, p *permission.Permission) error {
				p.ID = string(rune(i + 1))
				return nil
			}
			
			err := mockRepo.Create(ctx, &perms[i])
			assert.NoError(t, err)
			assert.NotEmpty(t, perms[i].ID)
		}
	})
}
