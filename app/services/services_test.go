package services_test

import (
	"GoHub-Service/app/services"
	"testing"
)

func TestTopicServiceCreate(t *testing.T) {
	// 准备测试数据
	dto := services.TopicCreateDTO{
		Title:      "测试话题",
		Body:       "这是一个测试话题的内容",
		CategoryID: "1",
		UserID:     "1",
	}

	// 创建服务实例
	service := services.NewTopicService()

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
	service := services.NewTopicService()

	// 测试获取不存在的话题
	_, err := service.GetByID("99999")
	if err == nil {
		t.Error("期望返回错误，但没有返回")
	}
}

func TestCategoryServiceCreate(t *testing.T) {
	dto := services.CategoryCreateDTO{
		Name:        "测试分类",
		Description: "这是一个测试分类",
	}

	service := services.NewCategoryService()
	category, err := service.Create(dto)

	if err != nil {
		t.Errorf("创建分类失败: %v", err)
		return
	}

	if category.Name != dto.Name {
		t.Errorf("分类名称不匹配, 期望: %s, 实际: %s", dto.Name, category.Name)
	}
}

func TestUserServiceList(t *testing.T) {
	// 注意: 这个测试需要 gin.Context，实际项目中应该使用 mock
	// 这里仅作为示例
	service := services.NewUserService()
	if service == nil {
		t.Error("服务创建失败")
	}
}
