// Package examples 代码重复消除使用示例
package examples

import (
	"context"
	"time"

	"GoHub-Service/app/models/comment"
	"GoHub-Service/app/models/topic"
	"GoHub-Service/pkg/mapper"
	"GoHub-Service/pkg/paginator"
	"GoHub-Service/pkg/repository"

	"github.com/gin-gonic/gin"
)

// ===== 示例 1: DTO 转换重复问题 =====

// ❌ 旧代码：每个 Service 都有重复的 toResponseDTO 方法

type OldCommentService struct{}

type OldCommentDTO struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// 重复代码 1
func (s *OldCommentService) toResponseDTO(c *comment.Comment) *OldCommentDTO {
	return &OldCommentDTO{
		ID:        c.GetStringID(),
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
	}
}

// 重复代码 2
func (s *OldCommentService) toResponseDTOList(comments []comment.Comment) []OldCommentDTO {
	dtos := make([]OldCommentDTO, len(comments))
	for i := range comments {
		dtos[i] = OldCommentDTO{
			ID:        comments[i].GetStringID(),
			Content:   comments[i].Content,
			CreatedAt: comments[i].CreatedAt,
		}
	}
	return dtos
}

// ✅ 新代码：使用通用 Mapper

type NewCommentService struct {
	mapper mapper.Mapper[comment.Comment, OldCommentDTO]
}

func NewNewCommentService() *NewCommentService {
	// 只需定义一次转换函数
	converter := func(c *comment.Comment) *OldCommentDTO {
		return &OldCommentDTO{
			ID:        c.GetStringID(),
			Content:   c.Content,
			CreatedAt: c.CreatedAt,
		}
	}

	return &NewCommentService{
		mapper: mapper.NewSimpleMapper(converter),
	}
}

// 不再需要 toResponseDTO 和 toResponseDTOList
func (s *NewCommentService) GetComment(c *comment.Comment) *OldCommentDTO {
	return s.mapper.ToDTO(c)
}

func (s *NewCommentService) GetComments(comments []comment.Comment) []OldCommentDTO {
	return s.mapper.ToDTOList(comments)
}

// ===== 示例 2: Repository CRUD 重复问题 =====

// ❌ 旧代码：每个 Repository 都有重复的 CRUD 方法

type OldTopicRepository struct{}

func (r *OldTopicRepository) GetByID(ctx context.Context, id string) (*topic.Topic, error) {
	// 重复的查询代码
	return nil, nil
}

func (r *OldTopicRepository) Create(ctx context.Context, t *topic.Topic) error {
	// 重复的创建代码
	return nil
}

func (r *OldTopicRepository) Update(ctx context.Context, t *topic.Topic) error {
	// 重复的更新代码
	return nil
}

func (r *OldTopicRepository) Delete(ctx context.Context, id string) error {
	// 重复的删除代码
	return nil
}

// ✅ 新代码：使用泛型 Repository 基类

type NewTopicRepository struct {
	*repository.GenericRepository[topic.Topic]
}

func NewNewTopicRepository() *NewTopicRepository {
	return &NewTopicRepository{
		GenericRepository: repository.NewGenericRepository[topic.Topic](),
	}
}

// 基础 CRUD 已由 GenericRepository 提供，无需重复编写
// 只需添加特定业务逻辑

func (r *NewTopicRepository) ListByCategory(ctx context.Context, c *gin.Context, categoryID string, perPage int) ([]topic.Topic, *paginator.Paging, error) {
	// 使用基类提供的 ListWithCondition
	return r.ListWithCondition(ctx, c, "category_id = ?", []interface{}{categoryID}, perPage)
}

func (r *NewTopicRepository) IncrementViewCount(ctx context.Context, id string) error {
	// 使用基类提供的 Increment
	return r.Increment(ctx, id, "view_count", 1)
}

// ===== 示例 3: 完整 Service 示例 =====

type ModernTopicService struct {
	repo   *NewTopicRepository
	mapper mapper.Mapper[topic.Topic, TopicDTO]
}

type TopicDTO struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	ViewCount int64     `json:"view_count"`
	CreatedAt time.Time `json:"created_at"`
}

func NewModernTopicService() *ModernTopicService {
	// 定义转换函数
	converter := func(t *topic.Topic) *TopicDTO {
		return &TopicDTO{
			ID:        t.GetStringID(),
			Title:     t.Title,
			ViewCount: t.ViewCount,
			CreatedAt: t.CreatedAt,
		}
	}

	return &ModernTopicService{
		repo:   NewNewTopicRepository(),
		mapper: mapper.NewSimpleMapper(converter),
	}
}

// 获取单个话题
func (s *ModernTopicService) GetTopic(ctx context.Context, id string) (*TopicDTO, error) {
	topicModel, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(topicModel), nil
}

// 获取话题列表
func (s *ModernTopicService) GetTopics(ctx context.Context, c *gin.Context, perPage int) ([]TopicDTO, *paginator.Paging, error) {
	topics, paging, err := s.repo.List(ctx, c, perPage)
	if err != nil {
		return nil, nil, err
	}
	return s.mapper.ToDTOList(topics), paging, nil
}

// 创建话题
func (s *ModernTopicService) CreateTopic(ctx context.Context, dto *TopicDTO) (*TopicDTO, error) {
	topicModel := &topic.Topic{
		Title:     dto.Title,
		ViewCount: 0,
	}
	
	if err := s.repo.Create(ctx, topicModel); err != nil {
		return nil, err
	}
	
	return s.mapper.ToDTO(topicModel), nil
}

// 增加浏览量
func (s *ModernTopicService) IncrementView(ctx context.Context, id string) error {
	return s.repo.IncrementViewCount(ctx, id)
}

// ===== 示例 4: 批量操作示例 =====

func (s *ModernTopicService) BatchCreateTopics(ctx context.Context, dtos []TopicDTO) error {
	// 转换 DTO 到 Model
	topics := make([]topic.Topic, len(dtos))
	for i, dto := range dtos {
		topics[i] = topic.Topic{
			Title: dto.Title,
		}
	}
	
	// 使用基类的批量创建（自动分块）
	return s.repo.BatchCreateInChunks(ctx, topics, 100)
}

func (s *ModernTopicService) BatchDeleteTopics(ctx context.Context, ids []string) error {
	// 使用基类的批量删除
	return s.repo.BatchDelete(ctx, ids)
}

// ===== 示例 5: 性能优化 - 大数据量并发转换 =====

type LargeDataService struct {
	repo   *NewTopicRepository
	mapper mapper.Mapper[topic.Topic, TopicDTO]
}

func NewLargeDataService() *LargeDataService {
	converter := func(t *topic.Topic) *TopicDTO {
		// 模拟复杂转换（如关联查询、计算等）
		time.Sleep(1 * time.Millisecond)
		return &TopicDTO{
			ID:        t.GetStringID(),
			Title:     t.Title,
			ViewCount: t.ViewCount,
			CreatedAt: t.CreatedAt,
		}
	}

	return &LargeDataService{
		repo:   NewNewTopicRepository(),
		mapper: mapper.NewBatchMapper(converter, 8), // 使用并发 Mapper
	}
}

func (s *LargeDataService) GetLargeList(ctx context.Context, c *gin.Context) ([]TopicDTO, error) {
	// 假设返回 1000+ 条数据
	topics, _, err := s.repo.List(ctx, c, 1000)
	if err != nil {
		return nil, err
	}
	
	// 并发转换（对于大数据量会自动启用）
	return s.mapper.ToDTOList(topics), nil
}

// ===== 示例 6: 代码行数对比 =====

/*
旧代码 (OldCommentService):
- toResponseDTO: 7 行
- toResponseDTOList: 11 行
- 总计: 18 行
- 每个 Service 都需要重复这 18 行

新代码 (NewCommentService):
- 转换函数: 6 行（只定义一次）
- 无需 toResponseDTO 和 toResponseDTOList
- 总计: 6 行
- **节省 66% 代码**

如果有 10 个 Service，节省：
- 旧代码: 18 × 10 = 180 行
- 新代码: 6 × 10 = 60 行
- **节省 120 行代码**

Repository 同理：
- 旧代码: 每个 Repository ~80 行基础 CRUD
- 新代码: 继承 GenericRepository，只需 ~10 行
- **节省 70 行/Repository**

如果有 10 个 Repository，节省：
- 旧代码: 80 × 10 = 800 行
- 新代码: 10 × 10 = 100 行
- **节省 700 行代码**

总计节省：120 + 700 = 820 行代码！
*/
