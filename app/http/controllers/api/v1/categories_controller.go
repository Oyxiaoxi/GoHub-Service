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

func (ctrl *CategoriesController) Store(c *gin.Context) {
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

func (ctrl *CategoriesController) Update(c *gin.Context) {
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
