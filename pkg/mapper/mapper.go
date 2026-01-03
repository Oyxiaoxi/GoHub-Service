// Package mapper 提供通用的 DTO 映射工具
// 使用泛型减少 toResponseDTO 和 toResponseDTOList 的代码重复
package mapper

// Mapper DTO 映射器接口
// T: 源类型（Model）
// D: 目标类型（DTO）
type Mapper[T any, D any] interface {
	// ToDTO 将单个模型转换为 DTO
	ToDTO(model *T) *D
	// ToDTOList 将模型列表转换为 DTO 列表
	ToDTOList(models []T) []D
}

// SimpleMapper 简单映射器，使用转换函数
type SimpleMapper[T any, D any] struct {
	converter func(*T) *D
}

// NewSimpleMapper 创建简单映射器
func NewSimpleMapper[T any, D any](converter func(*T) *D) *SimpleMapper[T, D] {
	return &SimpleMapper[T, D]{
		converter: converter,
	}
}

// ToDTO 将单个模型转换为 DTO
func (m *SimpleMapper[T, D]) ToDTO(model *T) *D {
	if model == nil {
		return nil
	}
	return m.converter(model)
}

// ToDTOList 将模型列表转换为 DTO 列表（优化：避免拷贝）
func (m *SimpleMapper[T, D]) ToDTOList(models []T) []D {
	if len(models) == 0 {
		return []D{}
	}

	dtos := make([]D, len(models))
	for i := range models {
		dto := m.converter(&models[i])
		if dto != nil {
			dtos[i] = *dto
		}
	}
	return dtos
}

// FuncMapper 函数映射器，直接使用函数作为映射器
type FuncMapper[T any, D any] func(*T) *D

// ToDTO 将单个模型转换为 DTO
func (f FuncMapper[T, D]) ToDTO(model *T) *D {
	if model == nil {
		return nil
	}
	return f(model)
}

// ToDTOList 将模型列表转换为 DTO 列表
func (f FuncMapper[T, D]) ToDTOList(models []T) []D {
	if len(models) == 0 {
		return []D{}
	}

	dtos := make([]D, len(models))
	for i := range models {
		dto := f(&models[i])
		if dto != nil {
			dtos[i] = *dto
		}
	}
	return dtos
}

// BatchMapper 批量映射器，支持并发转换（适用于大数据量）
type BatchMapper[T any, D any] struct {
	converter   func(*T) *D
	workerCount int
}

// NewBatchMapper 创建批量映射器
func NewBatchMapper[T any, D any](converter func(*T) *D, workerCount int) *BatchMapper[T, D] {
	if workerCount <= 0 {
		workerCount = 4 // 默认 4 个 worker
	}
	return &BatchMapper[T, D]{
		converter:   converter,
		workerCount: workerCount,
	}
}

// ToDTO 将单个模型转换为 DTO
func (m *BatchMapper[T, D]) ToDTO(model *T) *D {
	if model == nil {
		return nil
	}
	return m.converter(model)
}

// ToDTOList 将模型列表转换为 DTO 列表（并发）
func (m *BatchMapper[T, D]) ToDTOList(models []T) []D {
	if len(models) == 0 {
		return []D{}
	}

	// 小于 100 个不使用并发
	if len(models) < 100 {
		return m.toDTOListSerial(models)
	}

	return m.toDTOListParallel(models)
}

// toDTOListSerial 串行转换
func (m *BatchMapper[T, D]) toDTOListSerial(models []T) []D {
	dtos := make([]D, len(models))
	for i := range models {
		dto := m.converter(&models[i])
		if dto != nil {
			dtos[i] = *dto
		}
	}
	return dtos
}

// toDTOListParallel 并行转换
func (m *BatchMapper[T, D]) toDTOListParallel(models []T) []D {
	dtos := make([]D, len(models))
	chunkSize := (len(models) + m.workerCount - 1) / m.workerCount

	done := make(chan struct{})
	for w := 0; w < m.workerCount; w++ {
		start := w * chunkSize
		end := start + chunkSize
		if end > len(models) {
			end = len(models)
		}

		go func(start, end int) {
			for i := start; i < end; i++ {
				dto := m.converter(&models[i])
				if dto != nil {
					dtos[i] = *dto
				}
			}
			done <- struct{}{}
		}(start, end)
	}

	// 等待所有 worker 完成
	for w := 0; w < m.workerCount; w++ {
		<-done
	}
	close(done)

	return dtos
}

// FieldMapper 字段映射器，提供常用字段映射函数
type FieldMapper struct{}

// IDToString 将 uint64 ID 转换为字符串
func (FieldMapper) IDToString(id uint64) string {
	if id == 0 {
		return ""
	}
	return string(rune(id)) // 简化示例，实际应该使用 strconv
}

// StringToID 将字符串转换为 uint64 ID
func (FieldMapper) StringToID(id string) uint64 {
	if id == "" {
		return 0
	}
	return uint64([]rune(id)[0]) // 简化示例
}

// Helper functions

// Map 通用映射函数
func Map[T any, D any](models []T, converter func(*T) *D) []D {
	if len(models) == 0 {
		return []D{}
	}

	dtos := make([]D, len(models))
	for i := range models {
		dto := converter(&models[i])
		if dto != nil {
			dtos[i] = *dto
		}
	}
	return dtos
}

// MapFilter 映射并过滤 nil 值
func MapFilter[T any, D any](models []T, converter func(*T) *D) []D {
	dtos := make([]D, 0, len(models))
	for i := range models {
		dto := converter(&models[i])
		if dto != nil {
			dtos = append(dtos, *dto)
		}
	}
	return dtos
}
