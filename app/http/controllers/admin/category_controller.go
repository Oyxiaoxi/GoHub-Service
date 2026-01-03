package admin

import (
	"GoHub-Service/app/models/category"
	"GoHub-Service/pkg/database"
	"GoHub-Service/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CategoryController 分类管理控制器
type CategoryController struct{}

// Index 分类列表
func (ctrl *CategoryController) Index(c *gin.Context) {
	var categories []category.Category
	
	db := database.DB

	// 搜索条件
	if keyword := c.Query("keyword"); keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}

	db.Order("id ASC").Find(&categories)

	response.Data(c, gin.H{
		"categories": categories,
	})
}

// Show 分类详情
func (ctrl *CategoryController) Show(c *gin.Context) {
	categoryID := c.Param("id")

	var cat category.Category
	if err := database.DB.First(&cat, categoryID).Error; err != nil {
		response.Abort404(c, "分类不存在")
		return
	}

	// 统计该分类下的话题数
	var topicCount int64
	database.DB.Table("topics").Where("category_id = ?", categoryID).Count(&topicCount)

	response.Data(c, gin.H{
		"category":    cat,
		"topic_count": topicCount,
	})
}

// Store 创建分类
func (ctrl *CategoryController) Store(c *gin.Context) {
	type CreateCategoryRequest struct {
		Name        string `json:"name" binding:"required,min=2,max=50"`
		Description string `json:"description"`
	}

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, "参数错误")
		return
	}

	cat := category.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := database.DB.Create(&cat).Error; err != nil {
		response.Abort500(c, "创建失败")
		return
	}

	response.Created(c, gin.H{
		"category": cat,
	})
}

// Update 更新分类
func (ctrl *CategoryController) Update(c *gin.Context) {
	categoryID := c.Param("id")

	var cat category.Category
	if err := database.DB.First(&cat, categoryID).Error; err != nil {
		response.Abort404(c, "分类不存在")
		return
	}

	type UpdateCategoryRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, "参数错误")
		return
	}

	// 更新字段
	if req.Name != "" {
		cat.Name = req.Name
	}
	if req.Description != "" {
		cat.Description = req.Description
	}

	if err := database.DB.Save(&cat).Error; err != nil {
		response.Abort500(c, "更新失败")
		return
	}

	response.Data(c, gin.H{
		"category": cat,
	})
}

// Delete 删除分类
func (ctrl *CategoryController) Delete(c *gin.Context) {
	categoryID := c.Param("id")

	var cat category.Category
	if err := database.DB.First(&cat, categoryID).Error; err != nil {
		response.Abort404(c, "分类不存在")
		return
	}

	// 检查该分类下是否有话题
	var topicCount int64
	database.DB.Table("topics").Where("category_id = ?", categoryID).Count(&topicCount)
	
	if topicCount > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "该分类下还有话题，无法删除",
		})
		return
	}

	if err := database.DB.Delete(&cat).Error; err != nil {
		response.Abort500(c, "删除失败")
		return
	}

	response.Data(c, gin.H{
		"message": "分类已删除",
	})
}

// Sort 分类排序
func (ctrl *CategoryController) Sort(c *gin.Context) {
	type SortItem struct {
		ID    uint `json:"id" binding:"required"`
		Order int  `json:"order" binding:"required"`
	}
	
	type SortRequest struct {
		Items []SortItem `json:"items" binding:"required"`
	}

	var req SortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, "参数错误")
		return
	}

	// 批量更新排序
	for _, item := range req.Items {
		if err := database.DB.Model(&category.Category{}).
			Where("id = ?", item.ID).
			Update("sort_order", item.Order).Error; err != nil {
			response.Abort500(c, "排序失败")
			return
		}
	}

	response.Data(c, gin.H{
		"message": "排序成功",
		"count":   len(req.Items),
	})
}
