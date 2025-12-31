# GoHub-Service 测试指南

## 目录

- [1. 概述](#1-概述)
- [2. 测试策略](#2-测试策略)
- [3. 测试工具](#3-测试工具)
- [4. 单元测试规范](#4-单元测试规范)
- [5. Mock对象使用](#5-mock对象使用)
- [6. 测试覆盖率](#6-测试覆盖率)
- [7. 最佳实践](#7-最佳实践)
- [8. CI/CD集成](#8-cicd集成)

---

## 1. 概述

### 1.1 测试目标

- **覆盖率目标**: 60%+ (当前: ~40%)
- **核心层覆盖**: Service层 > 70%, Repository层 > 60%
- **关键路径**: 100%覆盖所有核心业务逻辑
- **持续改进**: 每个PR必须包含相应的测试用例

### 1.2 测试金字塔

```
       /\      E2E测试 (10%)
      /  \     集成测试 (30%)
     /    \    单元测试 (60%)
    /______\
```

## 2. 测试策略

### 2.1 测试层次

#### 单元测试 (Unit Tests)
- **范围**: 测试单个函数、方法
- **特点**: 快速、隔离、无依赖
- **工具**: Go testing, testutil
- **覆盖**: Service层、Repository层、工具函数

#### 集成测试 (Integration Tests)
- **范围**: 测试多个组件交互
- **特点**: 需要数据库、Redis等依赖
- **工具**: testcontainers, sqlmock
- **覆盖**: API端点、数据流程

#### E2E测试 (End-to-End Tests)
- **范围**: 测试完整的用户场景
- **特点**: 真实环境、完整流程
- **工具**: Postman, Newman
- **覆盖**: 核心业务流程

### 2.2 测试优先级

**P0 - 最高优先级（必须100%覆盖）**:
- 用户认证与授权
- 支付相关功能
- 数据安全功能（敏感词过滤、XSS防护）

**P1 - 高优先级（目标80%覆盖）**:
- 核心CRUD操作
- 业务逻辑Service层
- 数据访问Repository层

**P2 - 中优先级（目标60%覆盖）**:
- 辅助工具函数
- 中间件
- 缓存逻辑

**P3 - 低优先级（目标40%覆盖）**:
- 配置加载
- 日志记录
- 静态方法

## 3. 测试工具

### 3.1 内置工具包

项目提供了 `pkg/testutil` 包，包含常用测试辅助工具：

#### 断言助手 (assert_helper.go)

```go
import "GoHub-Service/pkg/testutil"

// 基础断言
testutil.AssertEqual(t, expected, actual, "值应该相等")
testutil.AssertNotEqual(t, expected, actual, "值应该不相等")
testutil.AssertNil(t, value, "应该为nil")
testutil.AssertNotNil(t, value, "不应该为nil")

// 布尔断言
testutil.AssertTrue(t, condition, "条件应该为true")
testutil.AssertFalse(t, condition, "条件应该为false")

// 错误断言
testutil.AssertNoError(t, err, "不应该有错误")
testutil.AssertError(t, err, "应该返回错误")

// 字符串断言
testutil.AssertContains(t, str, substr, "应该包含子串")

// 数值比较
testutil.AssertGreaterThan(t, a, b, "a应该大于b")
testutil.AssertLessThan(t, a, b, "a应该小于b")

// JSON比较
testutil.AssertJSONEqual(t, expectedJSON, actualJSON, "JSON应该相等")

// Panic断言
testutil.AssertPanic(t, func() { /* 会panic的代码 */ }, "应该panic")
```

#### 测试数据工厂 (mock_factory.go)

```go
import "GoHub-Service/pkg/testutil"

// 创建测试用户
user := testutil.MockUserFactory("1", "张三", "zhangsan@example.com")

// 创建测试分类
category := testutil.MockCategoryFactory("1", "技术", "技术讨论")

// 创建测试话题
topic := testutil.MockTopicFactory("1", "Go最佳实践", "内容...", "1", "1")

// 创建测试评论
comment := testutil.MockCommentFactory("1", "很好的内容", "1", "1", "")

// 批量创建
categories := testutil.MockCategories()  // 返回3个分类
topics := testutil.MockTopics()          // 返回3个话题
comments := testutil.MockComments()      // 返回3个评论

// 固定时间
mockTime := testutil.MockTime()
mockTimePtr := testutil.MockTimePtr()
```

### 3.2 第三方工具

推荐使用以下工具提升测试质量：

```bash
# Testify - 强大的断言库
go get github.com/stretchr/testify

# GoMock - Mock生成工具
go get github.com/golang/mock/gomock

# SQLMock - 数据库Mock
go get github.com/DATA-DOG/go-sqlmock

# HTTPMock - HTTP请求Mock
go get github.com/jarcoal/httpmock
```

## 4. 单元测试规范

### 4.1 测试文件命名

- 测试文件名: `xxx_test.go`
- 测试函数名: `TestXxx`
- 基准测试: `BenchmarkXxx`
- 示例函数: `ExampleXxx`

### 4.2 测试函数结构

#### Table-Driven Tests（推荐）

```go
func TestCategoryService_GetByID(t *testing.T) {
	tests := []struct {
		name      string           // 测试用例名称
		id        string           // 输入参数
		mockFunc  func(id string) (*category.Category, error)  // Mock行为
		wantErr   bool             // 是否期望错误
		checkFunc func(t *testing.T, result *CategoryResponseDTO)  // 结果检查
	}{
		{
			name: "成功获取分类",
			id:   "1",
			mockFunc: func(id string) (*category.Category, error) {
				return testutil.MockCategoryFactory("1", "技术", "技术讨论"), nil
			},
			wantErr: false,
			checkFunc: func(t *testing.T, result *CategoryResponseDTO) {
				testutil.AssertNotNil(t, result, "结果不应为nil")
				testutil.AssertEqual(t, "1", result.ID, "分类ID应该匹配")
				testutil.AssertEqual(t, "技术", result.Name, "分类名称应该匹配")
			},
		},
		{
			name: "分类不存在",
			id:   "999",
			mockFunc: func(id string) (*category.Category, error) {
				return nil, errors.New("分类不存在")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. 准备 - Setup
			mockRepo := &MockCategoryRepository{
				GetByIDFunc: tt.mockFunc,
			}
			service := &CategoryService{repo: mockRepo}

			// 2. 执行 - Execute
			result, err := service.GetByID(tt.id)

			// 3. 验证 - Verify
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
```

#### 3A模式 (Arrange-Act-Assert)

```go
func TestUserService_Create(t *testing.T) {
	// Arrange - 准备测试数据和依赖
	mockRepo := &MockUserRepository{
		CreateFunc: func(u *user.User) error {
			u.ID = "1"
			return nil
		},
	}
	service := &UserService{repo: mockRepo}
	dto := UserCreateDTO{
		Name:  "张三",
		Email: "zhangsan@example.com",
	}

	// Act - 执行被测试的方法
	result, err := service.Create(dto)

	// Assert - 断言结果
	testutil.AssertNil(t, err, "不应该有错误")
	testutil.AssertNotNil(t, result, "结果不应为nil")
	testutil.AssertEqual(t, "张三", result.Name, "用户名应该匹配")
}
```

### 4.3 测试覆盖的场景

每个函数至少应覆盖以下场景：

1. **正常场景（Happy Path）**
   - 输入有效数据
   - 期望成功返回

2. **边界场景（Edge Cases）**
   - 空值、零值
   - 最大值、最小值
   - 边界条件

3. **异常场景（Error Cases）**
   - 无效输入
   - 依赖失败
   - 业务规则违反

4. **并发场景（Concurrent）**（如果适用）
   - 多线程访问
   - 竞态条件

### 4.4 测试命名规范

测试用例名称应该清晰描述测试场景：

```go
// ✅ 好的命名
"成功创建分类"
"分类不存在时返回错误"
"空ID应该返回验证错误"
"并发创建多个分类"

// ❌ 不好的命名
"test1"
"error case"
"normal"
```

## 5. Mock对象使用

### 5.1 创建Mock Repository

```go
// MockCategoryRepository 分类仓储Mock
type MockCategoryRepository struct {
	GetByIDFunc func(id string) (*category.Category, error)
	ListFunc    func(c interface{}, perPage int) ([]category.Category, interface{}, error)
	CreateFunc  func(c *category.Category) error
	UpdateFunc  func(c *category.Category) error
	DeleteFunc  func(id string) error
}

func (m *MockCategoryRepository) GetByID(id string) (*category.Category, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	// 默认行为
	return testutil.MockCategoryFactory("1", "默认分类", "默认描述"), nil
}

// 确保实现了接口
var _ repositories.CategoryRepository = (*MockCategoryRepository)(nil)
```

### 5.2 使用Mock对象

```go
func TestService_WithMock(t *testing.T) {
	// 创建Mock，定义行为
	mockRepo := &MockCategoryRepository{
		GetByIDFunc: func(id string) (*category.Category, error) {
			if id == "999" {
				return nil, errors.New("不存在")
			}
			return testutil.MockCategoryFactory(id, "测试", "描述"), nil
		},
	}

	// 注入Mock
	service := &CategoryService{repo: mockRepo}

	// 测试
	result, err := service.GetByID("1")
	testutil.AssertNil(t, err, "不应该有错误")
	testutil.AssertEqual(t, "1", result.ID, "ID应该匹配")
}
```

### 5.3 Mock最佳实践

1. **每个接口一个Mock**: 为每个Repository接口创建独立的Mock
2. **默认行为**: 提供合理的默认返回值
3. **灵活配置**: 通过函数字段允许测试自定义行为
4. **类型安全**: 使用接口类型断言确保Mock实现了接口

## 6. 测试覆盖率

### 6.1 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./app/services/...

# 显示详细输出
go test -v ./...

# 运行特定测试
go test -run TestCategoryService_GetByID ./app/services/...
```

### 6.2 查看覆盖率

```bash
# 生成覆盖率报告
go test -cover ./...

# 生成详细覆盖率文件
go test -coverprofile=coverage.out ./...

# 查看覆盖率详情
go tool cover -func=coverage.out

# 生成HTML报告
go tool cover -html=coverage.out -o coverage.html

# 按包查看覆盖率
go test -cover ./app/services/...
go test -cover ./app/repositories/...
```

### 6.3 覆盖率目标

| 层次 | 当前 | 目标 | 说明 |
|-----|------|------|------|
| **总体** | ~40% | **60%+** | 整体项目覆盖率 |
| **Service层** | ~15% | **70%+** | 核心业务逻辑 |
| **Repository层** | ~10% | **60%+** | 数据访问层 |
| **工具包(pkg)** | ~5% | **50%+** | 公共工具函数 |
| **Controller层** | ~0% | **40%+** | HTTP处理层（需集成测试） |

### 6.4 提升覆盖率策略

**Phase 1 - 基础覆盖（当前 → 50%）**:
- ✅ 创建测试工具包（testutil）
- ✅ 为核心Service添加单元测试
- ✅ 为核心Repository添加单元测试
- ⏳ 为剩余Service添加测试
- ⏳ 为剩余Repository添加测试

**Phase 2 - 提升覆盖（50% → 60%）**:
- ⏳ pkg包单元测试
- ⏳ 中间件单元测试
- ⏳ 工具函数测试

**Phase 3 - 高覆盖（60% → 70%+）**:
- ⏳ Controller集成测试
- ⏳ 端到端测试
- ⏳ 性能测试

## 7. 最佳实践

### 7.1 测试编写原则

#### FIRST原则

- **Fast（快速）**: 单元测试应该在毫秒级完成
- **Independent（独立）**: 测试之间不应该有依赖
- **Repeatable（可重复）**: 任何环境都能重复执行
- **Self-Validating（自验证）**: 测试应该有明确的pass/fail结果
- **Timely（及时）**: 测试应该和代码一起编写

#### DRY原则

```go
// ❌ 不好的做法 - 重复代码
func TestCreate1(t *testing.T) {
	mockRepo := &MockRepo{...}
	service := &Service{repo: mockRepo}
	// ...
}

func TestCreate2(t *testing.T) {
	mockRepo := &MockRepo{...}
	service := &Service{repo: mockRepo}
	// ...
}

// ✅ 好的做法 - 提取公共逻辑
func setupService() (*Service, *MockRepo) {
	mockRepo := &MockRepo{...}
	service := &Service{repo: mockRepo}
	return service, mockRepo
}

func TestCreate1(t *testing.T) {
	service, _ := setupService()
	// ...
}

func TestCreate2(t *testing.T) {
	service, _ := setupService()
	// ...
}
```

### 7.2 常见陷阱

#### 1. 测试真实数据库

```go
// ❌ 不好 - 依赖真实数据库
func TestCreate(t *testing.T) {
	db := database.Connect()  // 真实连接
	repo := NewRepository(db)
	// ...
}

// ✅ 好 - 使用Mock
func TestCreate(t *testing.T) {
	mockRepo := &MockRepository{...}
	// ...
}
```

#### 2. 测试之间共享状态

```go
// ❌ 不好 - 共享状态
var sharedUser *user.User

func TestCreate(t *testing.T) {
	sharedUser = createUser()  // 影响其他测试
}

func TestUpdate(t *testing.T) {
	updateUser(sharedUser)  // 依赖TestCreate
}

// ✅ 好 - 每个测试独立
func TestCreate(t *testing.T) {
	user := testutil.MockUserFactory("1", "test", "test@example.com")
	// ...
}

func TestUpdate(t *testing.T) {
	user := testutil.MockUserFactory("1", "test", "test@example.com")
	// ...
}
```

#### 3. 过度Mock

```go
// ❌ 不好 - Mock过多
func TestComplexLogic(t *testing.T) {
	mockRepo1 := &MockRepo1{...}
	mockRepo2 := &MockRepo2{...}
	mockCache := &MockCache{...}
	mockLogger := &MockLogger{...}
	mockMetrics := &MockMetrics{...}
	// 太多Mock，测试变得脆弱
}

// ✅ 好 - 只Mock必要的依赖
func TestComplexLogic(t *testing.T) {
	mockRepo := &MockRepo{...}
	// 使用真实的Logger、Metrics（如果它们是简单的）
}
```

### 7.3 测试可维护性

#### 1. 使用辅助函数

```go
// 提取常用的测试设置
func setupCategoryTest() (*CategoryService, *MockCategoryRepository) {
	mockRepo := &MockCategoryRepository{}
	service := &CategoryService{repo: mockRepo}
	return service, mockRepo
}

// 提取常用的断言
func assertCategoryEqual(t *testing.T, expected, actual *CategoryResponseDTO) {
	testutil.AssertEqual(t, expected.ID, actual.ID, "ID应该匹配")
	testutil.AssertEqual(t, expected.Name, actual.Name, "名称应该匹配")
	testutil.AssertEqual(t, expected.Description, actual.Description, "描述应该匹配")
}
```

#### 2. 使用测试fixture

```go
// fixtures/category.go
package fixtures

func ValidCategory() *category.Category {
	return testutil.MockCategoryFactory("1", "技术", "技术讨论")
}

func InvalidCategory() *category.Category {
	return &category.Category{Name: ""}  // 无效分类
}

// 在测试中使用
func TestCreate(t *testing.T) {
	cat := fixtures.ValidCategory()
	// ...
}
```

### 7.4 测试文档化

```go
// TestCategoryService_Create 测试创建分类功能
//
// 测试场景:
// 1. 成功创建分类 - 验证返回的分类数据正确
// 2. 创建失败 - 数据库错误时应该返回错误
// 3. 验证失败 - 无效的分类名称应该返回验证错误
//
// 依赖:
// - MockCategoryRepository
//
// 注意事项:
// - 使用testutil.MockCategoryFactory创建测试数据
// - 不依赖真实数据库
func TestCategoryService_Create(t *testing.T) {
	// ...
}
```

## 8. CI/CD集成

### 8.1 GitHub Actions配置

```yaml
# .github/workflows/test.yml
name: Tests

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run tests
        run: go test -v -cover ./...
      
      - name: Generate coverage report
        run: go test -coverprofile=coverage.out ./...
      
      - name: Check coverage
        run: |
          coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          echo "Coverage: $coverage%"
          if (( $(echo "$coverage < 60" | bc -l) )); then
            echo "Coverage is below 60%"
            exit 1
          fi
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out
```

### 8.2 Pre-commit Hook

```bash
# .git/hooks/pre-commit
#!/bin/bash

echo "Running tests..."
go test ./... || exit 1

echo "Checking coverage..."
coverage=$(go test -cover ./... 2>&1 | grep -oP 'coverage: \K[0-9.]+' | tail -1)
if (( $(echo "$coverage < 60" | bc -l) )); then
    echo "Coverage is below 60% ($coverage%)"
    exit 1
fi

echo "Tests passed! Coverage: $coverage%"
```

### 8.3 Make命令

在 `Makefile` 中添加测试相关命令：

```makefile
# 运行所有测试
test:
	go test -v ./...

# 运行测试并生成覆盖率报告
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 运行特定包的测试
test-services:
	go test -v -cover ./app/services/...

test-repositories:
	go test -v -cover ./app/repositories/...

# 运行测试并检查覆盖率阈值
test-coverage-check:
	@coverage=$$(go test -cover ./... 2>&1 | grep -oP 'coverage: \K[0-9.]+' | tail -1); \
	echo "Total coverage: $$coverage%"; \
	if [ $$(echo "$$coverage < 60" | bc) -eq 1 ]; then \
		echo "❌ Coverage is below 60%"; \
		exit 1; \
	else \
		echo "✅ Coverage meets requirement (60%+)"; \
	fi

# 清理测试缓存
test-clean:
	go clean -testcache

# 运行基准测试
benchmark:
	go test -bench=. -benchmem ./...

.PHONY: test test-coverage test-services test-repositories test-coverage-check test-clean benchmark
```

使用方法：

```bash
make test                    # 运行所有测试
make test-coverage           # 生成覆盖率报告
make test-coverage-check     # 检查覆盖率是否达标
make test-services           # 只测试Service层
make test-repositories       # 只测试Repository层
```

## 9. 常见问题

### Q1: 如何测试需要数据库的代码？

**A**: 有三种方案：

1. **Mock Repository（推荐用于单元测试）**
   ```go
   mockRepo := &MockUserRepository{
       GetByIDFunc: func(id string) (*user.User, error) {
           return testutil.MockUserFactory(id, "test", "test@example.com"), nil
       },
   }
   ```

2. **SQLMock（用于测试SQL查询）**
   ```go
   db, mock, _ := sqlmock.New()
   mock.ExpectQuery("SELECT \\* FROM users").WillReturnRows(...)
   ```

3. **Testcontainers（用于集成测试）**
   ```go
   ctx := context.Background()
   mysqlContainer, _ := mysql.RunContainer(ctx, ...)
   ```

### Q2: 如何测试HTTP处理函数？

**A**: 使用 `httptest` 包：

```go
func TestCategoryController_Index(t *testing.T) {
	// 创建测试路由
	router := gin.Default()
	controller := NewCategoryController()
	router.GET("/categories", controller.Index)

	// 创建测试请求
	req := httptest.NewRequest("GET", "/categories", nil)
	w := httptest.NewRecorder()

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应
	testutil.AssertEqual(t, 200, w.Code, "状态码应该是200")
}
```

### Q3: 如何测试并发代码？

**A**: 使用 `sync.WaitGroup` 和并发安全断言：

```go
func TestConcurrentAccess(t *testing.T) {
	service := NewService()
	var wg sync.WaitGroup
	errors := make(chan error, 10)

	// 启动10个并发goroutine
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if err := service.DoSomething(id); err != nil {
				errors <- err
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// 检查是否有错误
	for err := range errors {
		t.Errorf("并发操作失败: %v", err)
	}
}
```

### Q4: 测试覆盖率低怎么办？

**A**: 按以下步骤提升：

1. **识别未覆盖的代码**
   ```bash
   go test -coverprofile=coverage.out ./...
   go tool cover -html=coverage.out
   ```

2. **优先覆盖核心逻辑**
   - Service层的业务逻辑
   - Repository层的数据操作
   - 工具函数

3. **使用测试模板**
   - 复用现有的测试模式
   - 使用testutil工具包

4. **持续改进**
   - 每个PR都要求增加测试
   - 定期Review测试覆盖率

## 10. 总结

### 10.1 关键要点

1. **测试优先**: 与代码一起编写测试
2. **保持简单**: 测试应该易于理解和维护
3. **快速反馈**: 单元测试应该在秒级完成
4. **真实场景**: 测试应该覆盖实际使用场景
5. **持续改进**: 定期Review和优化测试

### 10.2 检查清单

在提交PR前，确保：

- [ ] 所有新代码都有对应的测试
- [ ] 测试覆盖了正常、边界、异常场景
- [ ] 所有测试都能通过
- [ ] 测试覆盖率没有下降
- [ ] 使用了testutil工具包
- [ ] 测试有清晰的命名和注释
- [ ] 没有依赖外部服务（数据库、Redis等）
- [ ] 测试之间相互独立

### 10.3 参考资源

- [Go Testing官方文档](https://golang.org/pkg/testing/)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [Testify文档](https://github.com/stretchr/testify)
- [Go测试最佳实践](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)

---

**文档版本**: v1.0  
**更新日期**: 2025年12月31日  
**维护者**: GoHub Development Team

**下一步行动**:
1. 完成Service层剩余测试 → 提升至70%覆盖率
2. 完成Repository层剩余测试 → 提升至60%覆盖率
3. 添加集成测试 → API端点测试
4. 建立CI/CD测试流程 → 自动化质量检查

---

**© 2025 GoHub-Service Project. All rights reserved.**
