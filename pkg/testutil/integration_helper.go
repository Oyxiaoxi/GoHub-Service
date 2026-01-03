// Package integration 集成测试框架
package integration

import (
	"GoHub-Service/bootstrap"
	"GoHub-Service/config"
	"GoHub-Service/pkg/database"
	"GoHub-Service/pkg/logger"
	"context"
	"fmt"
	"os"
	"testing"

	"gorm.io/gorm"
)

// TestDB 测试数据库实例
var TestDB *gorm.DB

// TestConfig 测试配置
type TestConfig struct {
	DBPath      string
	LogLevel    string
	EnableCache bool
}

// DefaultTestConfig 返回默认测试配置
func DefaultTestConfig() *TestConfig {
	return &TestConfig{
		DBPath:      "./test_gohub.db",
		LogLevel:    "error",
		EnableCache: false,
	}
}

// SetupTestEnvironment 初始化测试环境
// 在测试开始前调用，初始化数据库、日志等
func SetupTestEnvironment(t *testing.T, cfg *TestConfig) {
	if cfg == nil {
		cfg = DefaultTestConfig()
	}

	// 设置测试环境变量
	os.Setenv("APP_ENV", "testing")
	os.Setenv("DB_CONNECTION", "sqlite")
	os.Setenv("DB_DATABASE", cfg.DBPath)

	// 初始化日志
	logger.InitForTest(cfg.LogLevel)

	// 初始化数据库
	db, err := setupTestDatabase(cfg.DBPath)
	if err != nil {
		t.Fatalf("初始化测试数据库失败: %v", err)
	}
	TestDB = db

	t.Cleanup(func() {
		TeardownTestEnvironment(t, cfg)
	})
}

// setupTestDatabase 初始化测试数据库
func setupTestDatabase(dbPath string) (*gorm.DB, error) {
	// 删除旧的测试数据库
	if _, err := os.Stat(dbPath); err == nil {
		os.Remove(dbPath)
	}

	// 创建新的测试数据库
	config.Initialize()
	db := database.Connect(dbPath, "sqlite")

	// 自动迁移表结构
	if err := autoMigrateTables(db); err != nil {
		return nil, fmt.Errorf("自动迁移表失败: %w", err)
	}

	return db, nil
}

// autoMigrateTables 自动迁移所有表
func autoMigrateTables(db *gorm.DB) error {
	// 这里需要导入所有模型并执行 AutoMigrate
	// 示例：
	// return db.AutoMigrate(
	//     &user.User{},
	//     &topic.Topic{},
	//     &category.Category{},
	//     &comment.Comment{},
	//     // ... 其他模型
	// )
	return nil
}

// TeardownTestEnvironment 清理测试环境
// 在测试结束后调用，清理数据库等资源
func TeardownTestEnvironment(t *testing.T, cfg *TestConfig) {
	if TestDB != nil {
		sqlDB, err := TestDB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	// 删除测试数据库文件
	if cfg != nil && cfg.DBPath != "" {
		os.Remove(cfg.DBPath)
	}
}

// BeginTransaction 开始事务
// 用于测试中的事务隔离
func BeginTransaction(t *testing.T) *gorm.DB {
	tx := TestDB.Begin()
	if tx.Error != nil {
		t.Fatalf("开始事务失败: %v", tx.Error)
	}

	t.Cleanup(func() {
		tx.Rollback()
	})

	return tx
}

// CleanTable 清空指定表
func CleanTable(db *gorm.DB, tableName string) error {
	return db.Exec(fmt.Sprintf("DELETE FROM %s", tableName)).Error
}

// CleanAllTables 清空所有测试表
func CleanAllTables(db *gorm.DB) error {
	tables := []string{
		"users",
		"topics",
		"categories",
		"comments",
		"links",
		"messages",
		"notifications",
		"roles",
		"permissions",
		"role_permissions",
		"user_roles",
		"likes",
		"follows",
	}

	for _, table := range tables {
		if err := CleanTable(db, table); err != nil {
			return fmt.Errorf("清空表 %s 失败: %w", table, err)
		}
	}

	return nil
}

// SeedTestData 填充测试数据
type SeedTestData func(db *gorm.DB) error

// RunWithTestData 使用测试数据运行测试
func RunWithTestData(t *testing.T, seed SeedTestData, testFunc func(db *gorm.DB)) {
	tx := BeginTransaction(t)

	// 填充测试数据
	if seed != nil {
		if err := seed(tx); err != nil {
			t.Fatalf("填充测试数据失败: %v", err)
		}
	}

	// 运行测试
	testFunc(tx)
}

// MockContext 创建测试用 Context
func MockContext() context.Context {
	ctx := context.Background()
	// 可以在这里添加 trace_id, user_id 等测试元数据
	return ctx
}

// AssertDBCount 断言表中记录数
func AssertDBCount(t *testing.T, db *gorm.DB, tableName string, expected int64) {
	var count int64
	if err := db.Table(tableName).Count(&count).Error; err != nil {
		t.Fatalf("统计表 %s 记录数失败: %v", tableName, err)
	}

	if count != expected {
		t.Errorf("表 %s 记录数不匹配: 期望 %d, 实际 %d", tableName, expected, count)
	}
}

// AssertRecordExists 断言记录存在
func AssertRecordExists(t *testing.T, db *gorm.DB, tableName string, condition string, args ...interface{}) {
	var count int64
	if err := db.Table(tableName).Where(condition, args...).Count(&count).Error; err != nil {
		t.Fatalf("查询表 %s 记录失败: %v", tableName, err)
	}

	if count == 0 {
		t.Errorf("表 %s 中不存在符合条件的记录: %s", tableName, condition)
	}
}

// TestHelper 集成测试辅助工具
type TestHelper struct {
	T  *testing.T
	DB *gorm.DB
}

// NewTestHelper 创建测试辅助工具
func NewTestHelper(t *testing.T, db *gorm.DB) *TestHelper {
	return &TestHelper{
		T:  t,
		DB: db,
	}
}

// AssertNoError 断言无错误
func (h *TestHelper) AssertNoError(err error, msg string) {
	if err != nil {
		h.T.Fatalf("%s: %v", msg, err)
	}
}

// AssertError 断言有错误
func (h *TestHelper) AssertError(err error, msg string) {
	if err == nil {
		h.T.Fatalf("%s: 期望有错误，但没有返回错误", msg)
	}
}

// AssertEqual 断言相等
func (h *TestHelper) AssertEqual(expected, actual interface{}, msg string) {
	if expected != actual {
		h.T.Errorf("%s: 期望 %v, 实际 %v", msg, expected, actual)
	}
}
