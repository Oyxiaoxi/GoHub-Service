// Package examples 资源管理使用示例
package examples

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"time"

	"GoHub-Service/pkg/logger"
	"GoHub-Service/pkg/resource"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ===== 示例 1: HTTP Response Body 安全关闭 =====

// ❌ 错误示例：可能泄漏
func BadHTTPRequest() error {
	resp, err := http.Get("https://api.example.com/data")
	if err != nil {
		return err
	}
	// 忘记关闭 Body，导致连接泄漏
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	return nil
}

// ✅ 正确示例：使用 defer 关闭
func GoodHTTPRequest() error {
	resp, err := http.Get("https://api.example.com/data")
	if err != nil {
		return err
	}
	defer resp.Body.Close() // 确保关闭
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	
	fmt.Println(string(body))
	return nil
}

// ✅ 更好示例：使用 SafeClose，捕获 panic
func BetterHTTPRequest() error {
	resp, err := http.Get("https://api.example.com/data")
	if err != nil {
		return err
	}
	defer resource.SafeClose(resp.Body, logger.Logger)
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	
	fmt.Println(string(body))
	return nil
}

// ===== 示例 2: 数据库事务安全处理 =====

// ❌ 错误示例：事务可能泄漏（GORM 已自动处理，这里是原始 sql.DB 示例）
func BadTransaction(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	
	// 如果这里 panic 或提前 return，事务不会回滚
	_, err = tx.Exec("INSERT INTO users (name) VALUES (?)", "alice")
	if err != nil {
		return err // 泄漏：未回滚
	}
	
	return tx.Commit()
}

// ✅ 正确示例：使用 defer 回滚
func GoodTransaction(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()
	
	_, err = tx.Exec("INSERT INTO users (name) VALUES (?)", "alice")
	if err != nil {
		return err
	}
	
	return tx.Commit()
}

// ✅ 更好示例：使用 TransactionGuard
func BetterTransaction(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	
	guard := resource.NewTransactionGuard(tx, logger.Logger)
	defer guard.Release()
	
	_, err = tx.Exec("INSERT INTO users (name) VALUES (?)", "alice")
	if err != nil {
		return err
	}
	
	if err := tx.Commit(); err != nil {
		return err
	}
	
	guard.Commit()
	return nil
}

// ✅ GORM 推荐示例：使用 Transaction 方法（自动处理）
func RecommendedGORMTransaction(db *gorm.DB) error {
	// GORM 的 Transaction 方法会自动处理 Commit 和 Rollback
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&User{Name: "alice"}).Error; err != nil {
			return err // 自动回滚
		}
		
		if err := tx.Create(&User{Name: "bob"}).Error; err != nil {
			return err // 自动回滚
		}
		
		return nil // 自动提交
	})
}

// ===== 示例 3: Context 取消防护 =====

// ❌ 错误示例：context 可能泄漏
func BadContextUsage() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 忘记调用 cancel，导致 goroutine 泄漏
	
	// 模拟耗时操作
	time.Sleep(1 * time.Second)
	fmt.Println("Done")
}

// ✅ 正确示例：使用 defer cancel
func GoodContextUsage() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // 确保释放资源
	
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("Done")
	case <-ctx.Done():
		fmt.Println("Timeout")
	}
}

// ✅ 更好示例：使用 ContextGuard
func BetterContextUsage() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	guard := resource.NewContextGuard(ctx, cancel, logger.Logger)
	defer guard.Release()
	
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("Done")
		guard.Cancel() // 手动取消
	case <-ctx.Done():
		fmt.Println("Timeout")
	}
}

// ===== 示例 4: Goroutine 池防止泄漏 =====

// ❌ 错误示例：无限制创建 goroutine
func BadGoroutineUsage() {
	for i := 0; i < 10000; i++ {
		go func(id int) {
			// 耗时操作
			time.Sleep(10 * time.Second)
			fmt.Printf("Task %d done\n", id)
		}(i)
		// 可能创建过多 goroutine，导致内存耗尽
	}
}

// ✅ 正确示例：使用 goroutine 池
func GoodGoroutineUsage() error {
	pool := resource.NewGoRoutinePool(10, logger.Logger) // 最多 10 个并发
	defer pool.Shutdown(30 * time.Second)
	
	for i := 0; i < 10000; i++ {
		taskID := i
		if err := pool.Submit(func() {
			// 耗时操作
			time.Sleep(1 * time.Second)
			fmt.Printf("Task %d done\n", taskID)
		}); err != nil {
			return err
		}
	}
	
	return nil
}

// ===== 示例 5: 资源泄漏检测 =====

// 使用 Tracker 追踪资源
func ResourceTrackingExample() {
	tracker := resource.NewTracker(logger.Logger)
	
	// 启动定期检查
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		
		for range ticker.C {
			tracker.Report(5 * time.Minute) // 报告超过 5 分钟未释放的资源
		}
	}()
	
	// 使用示例
	resp, err := http.Get("https://api.example.com/data")
	if err != nil {
		return
	}
	
	// 追踪资源
	resourceID := fmt.Sprintf("http-resp-%p", resp)
	tracker.Track(resourceID, "http.Response")
	defer func() {
		resp.Body.Close()
		tracker.Untrack(resourceID)
	}()
	
	// 使用资源...
}

// ===== 示例 6: Service 层完整示例 =====

type UserService struct {
	db      *gorm.DB
	pool    *resource.GoRoutinePool
	tracker *resource.Tracker
	logger  *zap.Logger
}

func NewUserService(db *gorm.DB, logger *zap.Logger) *UserService {
	return &UserService{
		db:      db,
		pool:    resource.NewGoRoutinePool(20, logger),
		tracker: resource.NewTracker(logger),
		logger:  logger,
	}
}

// CreateUser 创建用户（使用 GORM Transaction，自动处理回滚）
func (s *UserService) CreateUser(ctx context.Context, user *User) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建用户
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		
		// 创建用户设置
		settings := &UserSettings{UserID: user.ID}
		if err := tx.Create(settings).Error; err != nil {
			return err // 自动回滚
		}
		
		return nil // 自动提交
	})
}

// SendWelcomeEmail 发送欢迎邮件（使用 goroutine 池）
func (s *UserService) SendWelcomeEmail(ctx context.Context, userID int64) error {
	return s.pool.Submit(func() {
		// 异步发送邮件
		if err := s.sendEmail(ctx, userID); err != nil {
			s.logger.Error("发送欢迎邮件失败",
				zap.Int64("user_id", userID),
				zap.Error(err),
			)
		}
	})
}

// FetchUserData 获取用户数据（使用 context 超时控制）
func (s *UserService) FetchUserData(ctx context.Context, userID int64) (*User, error) {
	// 创建带超时的 context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	guard := resource.NewContextGuard(ctx, cancel, s.logger)
	defer guard.Release()
	
	var user User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return nil, err
	}
	
	guard.Cancel() // 手动取消（可选）
	return &user, nil
}

// Shutdown 关闭服务
func (s *UserService) Shutdown(timeout time.Duration) error {
	// 关闭 goroutine 池
	if err := s.pool.Shutdown(timeout); err != nil {
		s.logger.Error("关闭 goroutine 池失败", zap.Error(err))
		return err
	}
	
	// 检查资源泄漏
	s.tracker.Report(1 * time.Minute)
	
	return nil
}

// 辅助方法
func (s *UserService) sendEmail(ctx context.Context, userID int64) error {
	// 模拟发送邮件
	time.Sleep(1 * time.Second)
	return nil
}

// 数据模型
type User struct {
	ID   int64
	Name string
}

type UserSettings struct {
	ID     int64
	UserID int64
}
