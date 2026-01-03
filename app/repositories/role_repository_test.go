package repositories

import (
	"context"
	"testing"

	"github.com/Oyxiaoxi/GoHub-Service/app/models/role"
	"github.com/Oyxiaoxi/GoHub-Service/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupRoleRepositoryTest(t *testing.T) (*RoleRepository, *gorm.DB, func()) {
	helper := testutil.SetupTestEnvironment(t)
	repo := NewRoleRepository()
	return repo, helper.DB, helper.Close
}

func TestRoleRepository_Create(t *testing.T) {
	repo, db, cleanup := setupRoleRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("创建角色成功", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		testRole := testutil.MockRoleFactory()
		err := repo.Create(ctx, tx, testRole)

		assert.NoError(t, err)
		assert.NotZero(t, testRole.ID)
		assert.NotZero(t, testRole.CreatedAt)
	})

	t.Run("创建重复角色名称失败", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		testRole1 := testutil.MockRoleFactory()
		err := repo.Create(ctx, tx, testRole1)
		assert.NoError(t, err)

		// 相同角色名
		testRole2 := testutil.MockRoleFactory()
		testRole2.Name = testRole1.Name
		err = repo.Create(ctx, tx, testRole2)
		assert.Error(t, err)
	})
}

func TestRoleRepository_GetByID(t *testing.T) {
	repo, db, cleanup := setupRoleRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("获取存在的角色", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		testRole := testutil.MockRoleFactory()
		_ = repo.Create(ctx, tx, testRole)

		result, err := repo.GetByID(ctx, tx, testRole.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, testRole.ID, result.ID)
		assert.Equal(t, testRole.Name, result.Name)
	})

	t.Run("获取不存在的角色", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		result, err := repo.GetByID(ctx, tx, 99999)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestRoleRepository_GetByName(t *testing.T) {
	repo, db, cleanup := setupRoleRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("通过名称获取角色", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		testRole := testutil.MockRoleFactory()
		_ = repo.Create(ctx, tx, testRole)

		result, err := repo.GetByName(ctx, tx, testRole.Name)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, testRole.Name, result.Name)
	})

	t.Run("获取不存在的角色名", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		result, err := repo.GetByName(ctx, tx, "nonexistent-role")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestRoleRepository_GetAll(t *testing.T) {
	repo, db, cleanup := setupRoleRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("获取所有角色", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		// 创建测试角色
		roles := testutil.MockRoles(5)
		for _, r := range roles {
			_ = repo.Create(ctx, tx, r)
		}

		result, err := repo.GetAll(ctx, tx)

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(result), 5)
	})
}

func TestRoleRepository_Update(t *testing.T) {
	repo, db, cleanup := setupRoleRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("更新角色成功", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		testRole := testutil.MockRoleFactory()
		_ = repo.Create(ctx, tx, testRole)

		testRole.Description = "Updated Description"
		err := repo.Update(ctx, tx, testRole)

		assert.NoError(t, err)

		// 验证更新
		result, _ := repo.GetByID(ctx, tx, testRole.ID)
		assert.Equal(t, "Updated Description", result.Description)
	})
}

func TestRoleRepository_Delete(t *testing.T) {
	repo, db, cleanup := setupRoleRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("删除角色", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		testRole := testutil.MockRoleFactory()
		_ = repo.Create(ctx, tx, testRole)

		err := repo.Delete(ctx, tx, testRole.ID)

		assert.NoError(t, err)

		// 验证删除
		result, err := repo.GetByID(ctx, tx, testRole.ID)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestRoleRepository_Count(t *testing.T) {
	repo, db, cleanup := setupRoleRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("统计角色数量", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		// 创建测试角色
		roles := testutil.MockRoles(3)
		for _, r := range roles {
			_ = repo.Create(ctx, tx, r)
		}

		count, err := repo.Count(ctx, tx)

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(3))
	})
}

func TestRoleRepository_Paginate(t *testing.T) {
	repo, db, cleanup := setupRoleRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("分页获取角色", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		// 创建测试角色
		roles := testutil.MockRoles(10)
		for _, r := range roles {
			_ = repo.Create(ctx, tx, r)
		}

		result, total, err := repo.Paginate(ctx, tx, 1, 5)

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(10))
		assert.LessOrEqual(t, len(result), 5)
	})
}

func TestRoleRepository_Exists(t *testing.T) {
	repo, db, cleanup := setupRoleRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("检查角色是否存在", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		testRole := testutil.MockRoleFactory()
		_ = repo.Create(ctx, tx, testRole)

		exists := repo.Exists(ctx, tx, testRole.ID)
		assert.True(t, exists)

		notExists := repo.Exists(ctx, tx, 99999)
		assert.False(t, notExists)
	})
}

func TestRoleRepository_ExistsByName(t *testing.T) {
	repo, db, cleanup := setupRoleRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("通过名称检查角色是否存在", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		testRole := testutil.MockRoleFactory()
		_ = repo.Create(ctx, tx, testRole)

		exists := repo.ExistsByName(ctx, tx, testRole.Name)
		assert.True(t, exists)

		notExists := repo.ExistsByName(ctx, tx, "nonexistent-role")
		assert.False(t, notExists)
	})
}
