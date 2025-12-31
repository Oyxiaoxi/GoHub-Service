package auth

import (
    v1 "GoHub-Service/app/http/controllers/api/v1"
    "GoHub-Service/app/http/middlewares"
    "GoHub-Service/app/requests"
    "GoHub-Service/pkg/auth"
    "GoHub-Service/pkg/jwt"
    "GoHub-Service/pkg/logger"
    "GoHub-Service/pkg/response"

    "github.com/gin-gonic/gin"
)

// LoginController 用户控制器
type LoginController struct {
    v1.BaseAPIController
}

// LoginByPhone 手机登录
func (lc *LoginController) LoginByPhone(c *gin.Context) {

    // 1. 验证表单
    request := requests.LoginByPhoneRequest{}
    if ok := requests.Validate(c, &request, requests.LoginByPhone); !ok {
        return
    }

    // 2. 尝试登录
    user, err := auth.LoginByPhone(request.Phone)
    if err != nil {
        // 失败，显示错误提示
        response.Error(c, err, "账号不存在")
    } else {
        // 登录成功
        token := jwt.NewJWT().IssueToken(user.GetStringID(), user.Name)

        // 获取用户角色
        roles := middlewares.GetUserRoles(user.GetStringID())
        roleNames := make([]string, len(roles))
        for i, r := range roles {
            roleNames[i] = r.Name
        }

        // 返回用户信息和 token
        response.JSON(c, gin.H{
            "access_token": token,
            "user": gin.H{
                "id": user.GetStringID(),
                "name": user.Name,
                "email": user.Email,
                "phone": user.Phone,
                "avatar": user.Avatar,
                "city": user.City,
                "introduction": user.Introduction,
                "roles": roleNames,
                "created_at": user.CreatedAt,
                "updated_at": user.UpdatedAt,
            },
        })
    }
}

// LoginByPassword 多种方法登录，支持手机号、email 和用户名
func (lc *LoginController) LoginByPassword(c *gin.Context) {
    // 1. 验证表单
    request := requests.LoginByPasswordRequest{}
    if ok := requests.Validate(c, &request, requests.LoginByPassword); !ok {
        return
    }

    // 2. 尝试登录
    user, err := auth.Attempt(request.LoginID, request.Password)
    if err != nil {
        // 失败，显示错误提示
        logger.LogIf(err)
        response.Unauthorized(c, "账号不存在或密码错误")

    } else {
        token := jwt.NewJWT().IssueToken(user.GetStringID(), user.Name)
        
        // 获取用户角色
        roles := middlewares.GetUserRoles(user.GetStringID())
        roleNames := make([]string, len(roles))
        for i, r := range roles {
            roleNames[i] = r.Name
        }
        
        // 记录登陆日志
        logger.InfoString("Auth", "user_login", user.Email + " logged in successfully")
        
        // 返回用户信息和 token
        response.JSON(c, gin.H{
            "access_token": token,
            "user": gin.H{
                "id": user.GetStringID(),
                "name": user.Name,
                "email": user.Email,
                "phone": user.Phone,
                "avatar": user.Avatar,
                "city": user.City,
                "introduction": user.Introduction,
                "roles": roleNames,
                "created_at": user.CreatedAt,
                "updated_at": user.UpdatedAt,
            },
        })
    }
}

// RefreshToken 刷新 Access Token
func (lc *LoginController) RefreshToken(c *gin.Context) {

    token, err := jwt.NewJWT().RefreshToken(c)

    if err != nil {
        response.Error(c, err, "令牌刷新失败")
    } else {
        response.JSON(c, gin.H{
            "token": token,
        })
    }
}
