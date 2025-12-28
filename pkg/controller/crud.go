// Package controller 通用CRUD操作助手
package controller

import (
	"GoHub-Service/pkg/response"

	"github.com/gin-gonic/gin"
)

// Model 通用模型接口
type Model interface {
	GetID() uint64
	Create()
	Save() int64
	Delete() int64
}

// Validator 验证器接口
type Validator interface {
	Validate(*gin.Context) bool
}

// CRUDHelper CRUD操作助手结构
type CRUDHelper struct {
	ModelName string // 模型名称，用于错误提示
}

// NewCRUDHelper 创建CRUD助手
func NewCRUDHelper(modelName string) *CRUDHelper {
	return &CRUDHelper{
		ModelName: modelName,
	}
}

// HandleShow 处理Show操作
func (h *CRUDHelper) HandleShow(c *gin.Context, model Model) {
	if model.GetID() == 0 {
		response.Abort404(c)
		return
	}
	response.Data(c, model)
}

// HandleStore 处理Store操作
// modelFactory: 创建模型实例的函数
func (h *CRUDHelper) HandleStore(c *gin.Context, model Model) {
	model.Create()
	if model.GetID() > 0 {
		response.Created(c, model)
	} else {
		response.Abort500(c, h.ModelName+"创建失败，请稍后尝试~")
	}
}

// HandleUpdate 处理Update操作
func (h *CRUDHelper) HandleUpdate(c *gin.Context, model Model) {
	if model.GetID() == 0 {
		response.Abort404(c)
		return
	}

	rowsAffected := model.Save()
	if rowsAffected > 0 {
		response.Data(c, model)
	} else {
		response.Abort500(c, h.ModelName+"更新失败，请稍后尝试~")
	}
}

// HandleDelete 处理Delete操作
func (h *CRUDHelper) HandleDelete(c *gin.Context, model Model) {
	if model.GetID() == 0 {
		response.Abort404(c)
		return
	}

	rowsAffected := model.Delete()
	if rowsAffected > 0 {
		response.Success(c)
	} else {
		response.Abort500(c, h.ModelName+"删除失败，请稍后尝试~")
	}
}

// HandleList 处理列表查询
func (h *CRUDHelper) HandleList(c *gin.Context, data interface{}, pager interface{}) {
	response.JSON(c, gin.H{
		"data":  data,
		"pager": pager,
	})
}
