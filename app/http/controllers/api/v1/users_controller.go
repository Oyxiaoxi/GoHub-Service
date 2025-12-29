package v1

import (
    "GoHub-Service/app/requests"
    "GoHub-Service/app/services"
    "GoHub-Service/pkg/auth"
    "GoHub-Service/pkg/logger"
    "GoHub-Service/pkg/response"

    "github.com/gin-gonic/gin"
)

type UsersController struct {
    BaseAPIController
    userService *services.UserService
}

// NewUsersController 创建UsersController实例
func NewUsersController() *UsersController {
    return &UsersController{
        userService: services.NewUserService(),
    }
}

// CurrentUser 当前登录用户信息
func (ctrl *UsersController) CurrentUser(c *gin.Context) {
    userModel := auth.CurrentUser(c)
    response.Data(c, userModel)
}

// Index 所有用户
func (ctrl *UsersController) Index(c *gin.Context) {
    request := requests.PaginationRequest{}
    if ok := requests.Validate(c, &request, requests.Pagination); !ok {
        return
    }

    listResponse, err := ctrl.userService.List(c, 10)
    if err != nil {
        logger.LogErrorWithContext(c, err, "获取用户列表失败")
        response.Abort500(c, "获取列表失败")
        return
    }

    response.JSON(c, gin.H{
        "data":  listResponse.Users,
        "pager": listResponse.Paging,
    })
}

func (ctrl *UsersController) UpdateProfile(c *gin.Context) {

    request := requests.UserUpdateProfileRequest{}
    if ok := requests.Validate(c, &request, requests.UserUpdateProfile); !ok {
        return
    }

    currentUser := auth.CurrentUser(c)
    updatedUser, err := ctrl.userService.UpdateProfile(&currentUser, request.Name, request.City, request.Introduction)
    if err != nil {
        logger.LogErrorWithContext(c, err, "用户信息更新失败")
        response.Abort500(c, "更新失败，请稍后尝试~")
        return
    }
    response.Data(c, updatedUser)
}
