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

func (ctrl *LinksController) Index(c *gin.Context) {
    response.Data(c, ctrl.linkService.GetAllCached())
}
