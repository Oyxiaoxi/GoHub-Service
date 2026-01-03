package v1

import (
    "GoHub-Service/app/http/middlewares"
    "GoHub-Service/app/requests"
    "GoHub-Service/app/services"
    apperrors "GoHub-Service/pkg/errors"
    "GoHub-Service/pkg/logger"
    "GoHub-Service/pkg/response"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type CategoriesController struct {
    BaseAPIController
    categoryService *services.CategoryService
}

// NewCategoriesController 创建CategoriesController实例
func NewCategoriesController() *CategoriesController {
    return &CategoriesController{
        categoryService: services.NewCategoryService(),
    }
}
// Store 创建分类
// @Summary 创建新分类
// @Description 创建一个新的话题分类
// @Tags 分类管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param category body requests.CategoryRequest true "分类信息"
// @Success 201 {object} response.Response "成功"
// @Failure 422 {object} response.Response "验证失败"
// @Router /categories [post]func (ctrl *CategoriesController) Store(c *gin.Context) {
    request := requests.CategoryRequest{}
    if ok := requests.Validate(c, &request, requests.CategorySave); !ok {
        return
    }

    dto := services.CategoryCreateDTO{
        Name:        request.Name,
        Description: request.Description,
    }

    categoryModel, err := ctrl.categoryService.Create(dto)
    if err != nil {
        logger.LogErrorWithContext(c, err, "创建分类失败")
        if appErr, ok := err.(*apperrors.AppError); ok {
            response.ApiError(c, 500, appErr.Code, appErr.Message)
        } else {
            response.Abort500(c, "创建失败，请稍后尝试~")
        }
        return
    }

    response.Created(c, categoryModel)
}
// Update 更新分类
// @Summary 更新分类信息
// @Description 更新指定分类的信息
// @Tags 分类管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "分类ID"
// @Param category body requests.CategoryRequest true "分类信息"
// @Success 200 {object} response.Response "成功"
// @Failure 404 {object} response.Response "分类不存在"
// @Failure 422 {object} response.Response "验证失败"
// @Router /categories/{id} [put]func (ctrl *CategoriesController) Update(c *gin.Context) {
    // 表单验证
    request := requests.CategoryRequest{}
    if ok := requests.Validate(c, &request, requests.CategorySave); !ok {
        return
    }

    // 更新分类
    dto := services.CategoryUpdateDTO{
        Name:        &request.Name,
        Description: &request.Description,
    }

    categoryModel, err := ctrl.categoryService.Update(c.Param("id"), dto)
    if err != nil {
        if apperrors.IsAppError(err) {
            appErr, _ := apperrors.GetAppError(err)
            appErr.WithRequestID(middlewares.GetRequestID(c))
            response.Abort404(c)
            return
        }
        logger.LogErrorWithContext(c, err, "更新分类失败",
            zap.String("category_id", c.Param("id")),
        )
        response.Abort500(c, "更新失败，请稍后尝试~")
        return
    }

    response.Data(c, categoryModel)
}

func (ctrl *CategoriesController) Index(c *gin.Context) {
    request := requests.PaginationRequest{}
    if ok := requests.Validate(c, &request, requests.Pagination); !ok {
        return
    }

    listResponse, err := ctrl.categoryService.List(c, 10)
    if err != nil {
        logger.LogErrorWithContext(c, err, "获取分类列表失败")
        response.Abort500(c, "获取列表失败")
        return
    }

    response.JSON(c, gin.H{
        "data":  listResponse.Categories,
        "pager": listResponse.Paging,
    })
}

func (ctrl *CategoriesController) Delete(c *gin.Context) {
    err := ctrl.categoryService.Delete(c.Param("id"))
    if err != nil {
        if apperrors.IsAppError(err) {
            appErr, _ := apperrors.GetAppError(err)
            appErr.WithRequestID(middlewares.GetRequestID(c))
            response.Abort404(c)
            return
        }
        logger.LogErrorWithContext(c, err, "删除分类失败",
            zap.String("category_id", c.Param("id")),
        )
        response.Abort500(c, "删除失败，请稍后尝试~")
        return
    }

    response.Success(c)
}
