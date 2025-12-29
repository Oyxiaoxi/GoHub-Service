package v1

import (
    "GoHub-Service/app/requests"
    "GoHub-Service/app/services"
    "GoHub-Service/pkg/auth"
    "GoHub-Service/pkg/logger"
    "GoHub-Service/pkg/response"
    "GoHub-Service/pkg/file"
    "GoHub-Service/pkg/config"

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

func (ctrl *UsersController) UpdateEmail(c *gin.Context) {

    request := requests.UserUpdateEmailRequest{}
    if ok := requests.Validate(c, &request, requests.UserUpdateEmail); !ok {
        return
    }

    currentUser := auth.CurrentUser(c)
    currentUser.Email = request.Email
    rowsAffected := currentUser.Save()

    if rowsAffected > 0 {
        response.Success(c)
    } else {
        // 失败，显示错误提示
        response.Abort500(c, "更新失败，请稍后尝试~")
    }
}

func (ctrl *UsersController) UpdatePhone(c *gin.Context) {

    request := requests.UserUpdatePhoneRequest{}
    if ok := requests.Validate(c, &request, requests.UserUpdatePhone); !ok {
        return
    }

    currentUser := auth.CurrentUser(c)
    currentUser.Phone = request.Phone
    rowsAffected := currentUser.Save()

    if rowsAffected > 0 {
        response.Success(c)
    } else {
        response.Abort500(c, "更新失败，请稍后尝试~")
    }
}

func (ctrl *UsersController) UpdatePassword(c *gin.Context) {

    request := requests.UserUpdatePasswordRequest{}
    if ok := requests.Validate(c, &request, requests.UserUpdatePassword); !ok {
        return
    }

    currentUser := auth.CurrentUser(c)
    // 验证原始密码是否正确
    _, err := auth.Attempt(currentUser.Name, request.Password)
    if err != nil {
        // 失败，显示错误提示
        response.Unauthorized(c, "原密码不正确")
    } else {
        // 更新密码为新密码
        currentUser.Password = request.NewPassword
        currentUser.Save()

        response.Success(c)
    }
}

func (ctrl *UsersController) UpdateAvatar(c *gin.Context) {

    request := requests.UserUpdateAvatarRequest{}
    if ok := requests.Validate(c, &request, requests.UserUpdateAvatar); !ok {
        return
    }

    avatar, err := file.SaveUploadAvatar(c, request.Avatar)
    if err != nil {
        response.Abort500(c, "上传头像失败，请稍后尝试~")
        return
    }

    currentUser := auth.CurrentUser(c)
    currentUser.Avatar = config.GetString("app.url") + avatar
    currentUser.Save()

    response.Data(c, currentUser)
}
