// Package service 业务逻辑层基础定义
package service

// 说明：Service层负责处理业务逻辑，从Controller中分离
// 
// 优点：
// 1. 业务逻辑与HTTP逻辑分离
// 2. 便于单元测试
// 3. 代码复用性更好
// 4. 便于事务管理
//
// 使用示例：
//
// 1. 创建Service结构体：
//    type UserService struct {}
//
// 2. 实现业务方法：
//    func (s *UserService) CreateUser(data CreateUserDTO) (*user.User, error) {
//        // 业务逻辑处理
//        userModel := user.User{
//            Name: data.Name,
//            // ...
//        }
//        userModel.Create()
//        return &userModel, nil
//    }
//
// 3. 在Controller中调用：
//    func (ctrl *UsersController) Store(c *gin.Context) {
//        var request requests.UserRequest
//        if ok := requests.Validate(c, &request, requests.UserSave); !ok {
//            return
//        }
//        
//        userService := service.UserService{}
//        user, err := userService.CreateUser(request.ToDTO())
//        if err != nil {
//            response.ApiError(c, http.StatusInternalServerError, response.CodeCreateFailed, err.Error())
//            return
//        }
//        
//        response.ApiSuccess(c, user)
//    }

// BaseService 基础Service结构（可选）
type BaseService struct{}

// TransactionFunc 事务处理函数类型
type TransactionFunc func() error
