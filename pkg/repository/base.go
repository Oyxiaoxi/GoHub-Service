// Package repository 数据访问层基础定义
package repository

// 说明：Repository模式用于封装数据访问逻辑
//
// 优点：
// 1. 数据访问逻辑与业务逻辑分离
// 2. 便于切换数据源（MySQL -> PostgreSQL）
// 3. 统一数据访问接口
// 4. 便于Mock测试
//
// 使用示例：
//
// 1. 定义Repository接口：
//    type UserRepository interface {
//        FindByID(id uint64) (*user.User, error)
//        FindByEmail(email string) (*user.User, error)
//        Create(user *user.User) error
//        Update(user *user.User) error
//        Delete(id uint64) error
//        List(page, pageSize int) ([]user.User, int64, error)
//    }
//
// 2. 实现Repository：
//    type userRepository struct {
//        db *gorm.DB
//    }
//
//    func NewUserRepository(db *gorm.DB) UserRepository {
//        return &userRepository{db: db}
//    }
//
//    func (r *userRepository) FindByID(id uint64) (*user.User, error) {
//        var user user.User
//        err := r.db.Where("id = ?", id).First(&user).Error
//        return &user, err
//    }
//
// 3. 在Service中使用：
//    type UserService struct {
//        userRepo repository.UserRepository
//    }
//
//    func (s *UserService) GetUser(id uint64) (*user.User, error) {
//        return s.userRepo.FindByID(id)
//    }

// BaseRepository 基础Repository结构（可选）
type BaseRepository struct{}

// PaginationParams 分页参数
type PaginationParams struct {
	Page     int
	PageSize int
	OrderBy  string
}

// PaginationResult 分页结果
type PaginationResult struct {
	Total       int64       `json:"total"`
	Page        int         `json:"page"`
	PageSize    int         `json:"page_size"`
	TotalPages  int         `json:"total_pages"`
	Data        interface{} `json:"data"`
}
