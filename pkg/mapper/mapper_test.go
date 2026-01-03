package mapper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试模型
type TestModel struct {
	ID        uint64
	Name      string
	Email     string
	CreatedAt time.Time
}

// 测试 DTO
type TestDTO struct {
	ID        string
	Name      string
	Email     string
	CreatedAt string
}

// 转换函数
func modelToDTO(model *TestModel) *TestDTO {
	if model == nil {
		return nil
	}
	return &TestDTO{
		ID:        "id-123", // 简化示例
		Name:      model.Name,
		Email:     model.Email,
		CreatedAt: model.CreatedAt.Format("2006-01-02"),
	}
}

func TestSimpleMapper_ToDTO(t *testing.T) {
	mapper := NewSimpleMapper(modelToDTO)

	t.Run("正常转换", func(t *testing.T) {
		model := &TestModel{
			ID:        1,
			Name:      "Alice",
			Email:     "alice@example.com",
			CreatedAt: time.Now(),
		}

		dto := mapper.ToDTO(model)
		assert.NotNil(t, dto)
		assert.Equal(t, "Alice", dto.Name)
		assert.Equal(t, "alice@example.com", dto.Email)
	})

	t.Run("nil 输入", func(t *testing.T) {
		dto := mapper.ToDTO(nil)
		assert.Nil(t, dto)
	})
}

func TestSimpleMapper_ToDTOList(t *testing.T) {
	mapper := NewSimpleMapper(modelToDTO)

	t.Run("正常转换列表", func(t *testing.T) {
		models := []TestModel{
			{ID: 1, Name: "Alice", Email: "alice@example.com", CreatedAt: time.Now()},
			{ID: 2, Name: "Bob", Email: "bob@example.com", CreatedAt: time.Now()},
			{ID: 3, Name: "Charlie", Email: "charlie@example.com", CreatedAt: time.Now()},
		}

		dtos := mapper.ToDTOList(models)
		assert.Len(t, dtos, 3)
		assert.Equal(t, "Alice", dtos[0].Name)
		assert.Equal(t, "Bob", dtos[1].Name)
		assert.Equal(t, "Charlie", dtos[2].Name)
	})

	t.Run("空列表", func(t *testing.T) {
		dtos := mapper.ToDTOList([]TestModel{})
		assert.Empty(t, dtos)
	})
}

func TestFuncMapper(t *testing.T) {
	mapper := FuncMapper[TestModel, TestDTO](modelToDTO)

	t.Run("ToDTO", func(t *testing.T) {
		model := &TestModel{
			ID:        1,
			Name:      "Alice",
			Email:     "alice@example.com",
			CreatedAt: time.Now(),
		}

		dto := mapper.ToDTO(model)
		assert.NotNil(t, dto)
		assert.Equal(t, "Alice", dto.Name)
	})

	t.Run("ToDTOList", func(t *testing.T) {
		models := []TestModel{
			{ID: 1, Name: "Alice", Email: "alice@example.com", CreatedAt: time.Now()},
			{ID: 2, Name: "Bob", Email: "bob@example.com", CreatedAt: time.Now()},
		}

		dtos := mapper.ToDTOList(models)
		assert.Len(t, dtos, 2)
	})
}

func TestBatchMapper(t *testing.T) {
	mapper := NewBatchMapper(modelToDTO, 4)

	t.Run("小数据量（串行）", func(t *testing.T) {
		models := make([]TestModel, 50)
		for i := range models {
			models[i] = TestModel{
				ID:        uint64(i + 1),
				Name:      "User" + string(rune(i)),
				Email:     "user@example.com",
				CreatedAt: time.Now(),
			}
		}

		dtos := mapper.ToDTOList(models)
		assert.Len(t, dtos, 50)
	})

	t.Run("大数据量（并发）", func(t *testing.T) {
		models := make([]TestModel, 500)
		for i := range models {
			models[i] = TestModel{
				ID:        uint64(i + 1),
				Name:      "User" + string(rune(i)),
				Email:     "user@example.com",
				CreatedAt: time.Now(),
			}
		}

		dtos := mapper.ToDTOList(models)
		assert.Len(t, dtos, 500)
	})
}

func TestMap(t *testing.T) {
	models := []TestModel{
		{ID: 1, Name: "Alice", Email: "alice@example.com", CreatedAt: time.Now()},
		{ID: 2, Name: "Bob", Email: "bob@example.com", CreatedAt: time.Now()},
	}

	dtos := Map(models, modelToDTO)
	assert.Len(t, dtos, 2)
	assert.Equal(t, "Alice", dtos[0].Name)
	assert.Equal(t, "Bob", dtos[1].Name)
}

func TestMapFilter(t *testing.T) {
	// 模拟部分转换失败返回 nil
	converter := func(model *TestModel) *TestDTO {
		if model.Name == "" {
			return nil
		}
		return modelToDTO(model)
	}

	models := []TestModel{
		{ID: 1, Name: "Alice", Email: "alice@example.com", CreatedAt: time.Now()},
		{ID: 2, Name: "", Email: "empty@example.com", CreatedAt: time.Now()}, // 将被过滤
		{ID: 3, Name: "Charlie", Email: "charlie@example.com", CreatedAt: time.Now()},
	}

	dtos := MapFilter(models, converter)
	// 注意：当前实现不会过滤 nil，需要修改实现
	assert.Len(t, dtos, 2) // 预期过滤掉 1 个
}

// Benchmark 测试
func BenchmarkSimpleMapper(b *testing.B) {
	mapper := NewSimpleMapper(modelToDTO)
	models := make([]TestModel, 100)
	for i := range models {
		models[i] = TestModel{
			ID:        uint64(i + 1),
			Name:      "User",
			Email:     "user@example.com",
			CreatedAt: time.Now(),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mapper.ToDTOList(models)
	}
}

func BenchmarkBatchMapper_Serial(b *testing.B) {
	mapper := NewBatchMapper(modelToDTO, 4)
	models := make([]TestModel, 50) // 小于 100，使用串行
	for i := range models {
		models[i] = TestModel{
			ID:        uint64(i + 1),
			Name:      "User",
			Email:     "user@example.com",
			CreatedAt: time.Now(),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mapper.ToDTOList(models)
	}
}

func BenchmarkBatchMapper_Parallel(b *testing.B) {
	mapper := NewBatchMapper(modelToDTO, 4)
	models := make([]TestModel, 500) // 大于 100，使用并发
	for i := range models {
		models[i] = TestModel{
			ID:        uint64(i + 1),
			Name:      "User",
			Email:     "user@example.com",
			CreatedAt: time.Now(),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mapper.ToDTOList(models)
	}
}

func BenchmarkMap(b *testing.B) {
	models := make([]TestModel, 100)
	for i := range models {
		models[i] = TestModel{
			ID:        uint64(i + 1),
			Name:      "User",
			Email:     "user@example.com",
			CreatedAt: time.Now(),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Map(models, modelToDTO)
	}
}
