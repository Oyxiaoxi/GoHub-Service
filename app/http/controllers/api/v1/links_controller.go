package v1

import (
    "GoHub-Service/app/services"
    "GoHub-Service/pkg/response"

    "github.com/gin-gonic/gin"
)

type LinksController struct {
    BaseAPIController
    linkService *services.LinkService
}

// NewLinksController 创建LinksController实例
func NewLinksController() *LinksController {
    return &LinksController{
        linkService: services.NewLinkService(),
    }
}

// Index 友情链接列表
// @Summary 获取友情链接列表
// @Description 获取所有激活的友情链接，结果从缓存中获取
// @Tags 友情链接
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "成功"
// @Router /links [get]
func (ctrl *LinksController) Index(c *gin.Context) {
    listResponse, err := ctrl.linkService.GetAllCached()
    if err != nil {
        response.ApiError(c, 500, err.Code, err.Message)
        return
    }
    response.Data(c, listResponse)
}
