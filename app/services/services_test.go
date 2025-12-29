package services

import "testing"

func TestTopicServiceCreate(t *testing.T) {
	// 准备测试数据
	dto := TopicCreateDTO{
		Title:      "测试话题",
		Body:       "这是一个测试话题的内容",
		CategoryID: "1",
		UserID:     "1",
	}

	// 创建服务实例（使用 Mock，避免真实 DB）
	service := &TopicService{repo: NewMockTopicRepository()}

	// 执行创建操作
	topic, err := service.Create(dto)

	// 验证结果
	if err != nil {
		t.Errorf("创建话题失败: %v", err)
		return
	}

	if topic.Title != dto.Title {
		t.Errorf("话题标题不匹配, 期望: %s, 实际: %s", dto.Title, topic.Title)
	}

	if topic.Body != dto.Body {
		t.Errorf("话题内容不匹配, 期望: %s, 实际: %s", dto.Body, topic.Body)
	}
}

func TestTopicServiceGetByID(t *testing.T) {
	service := &TopicService{repo: NewMockTopicRepository()}

	// 测试获取不存在的话题
	_, err := service.GetByID("99999")
	if err == nil {
		t.Error("期望返回错误，但没有返回")
	}
}

func TestCategoryServiceCreate(t *testing.T) {
	t.Skip("Category service create requires real DB; skipped in unit test")
}

func TestUserServiceList(t *testing.T) {
	t.Skip("User service list requires gin context/DB; skipped in unit test")
}
