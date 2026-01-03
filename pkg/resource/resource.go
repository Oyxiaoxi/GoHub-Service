// Package resource 资源管理工具包
// 提供资源泄漏检测和防护功能
package resource

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Closer 可关闭的资源接口
type Closer interface {
	Close() error
}

// SafeClose 安全关闭资源，捕获 panic
func SafeClose(closer Closer, logger *zap.Logger) {
	if closer == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			if logger != nil {
				logger.Error("资源关闭时发生 panic",
					zap.Any("panic", r),
					zap.String("stack", string(debug.Stack())),
				)
			}
		}
	}()

	if err := closer.Close(); err != nil {
		if logger != nil {
			logger.Warn("资源关闭失败",
				zap.Error(err),
				zap.String("type", fmt.Sprintf("%T", closer)),
			)
		}
	}
}

// Tracker 资源追踪器，用于检测资源泄漏
type Tracker struct {
	mu        sync.RWMutex
	resources map[string]*ResourceInfo
	logger    *zap.Logger
}

// ResourceInfo 资源信息
type ResourceInfo struct {
	ID        string
	Type      string
	CreatedAt time.Time
	Stack     string
}

// NewTracker 创建资源追踪器
func NewTracker(logger *zap.Logger) *Tracker {
	return &Tracker{
		resources: make(map[string]*ResourceInfo),
		logger:    logger,
	}
}

// Track 追踪资源
func (t *Tracker) Track(id, resourceType string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.resources[id] = &ResourceInfo{
		ID:        id,
		Type:      resourceType,
		CreatedAt: time.Now(),
		Stack:     string(debug.Stack()),
	}
}

// Untrack 停止追踪资源
func (t *Tracker) Untrack(id string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.resources, id)
}

// Check 检查资源泄漏（超过指定时间未释放的资源）
func (t *Tracker) Check(threshold time.Duration) []*ResourceInfo {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var leaked []*ResourceInfo
	now := time.Now()

	for _, info := range t.resources {
		if now.Sub(info.CreatedAt) > threshold {
			leaked = append(leaked, info)
		}
	}

	return leaked
}

// Report 报告资源泄漏
func (t *Tracker) Report(threshold time.Duration) {
	leaked := t.Check(threshold)
	if len(leaked) == 0 {
		return
	}

	if t.logger != nil {
		for _, info := range leaked {
			t.logger.Warn("检测到可能的资源泄漏",
				zap.String("id", info.ID),
				zap.String("type", info.Type),
				zap.Time("created_at", info.CreatedAt),
				zap.Duration("age", time.Since(info.CreatedAt)),
				zap.String("stack", info.Stack),
			)
		}
	}
}

// Count 返回当前追踪的资源数量
func (t *Tracker) Count() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.resources)
}

// Clear 清除所有追踪信息（用于测试）
func (t *Tracker) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.resources = make(map[string]*ResourceInfo)
}

// GoRoutinePool goroutine 池，防止 goroutine 泄漏
type GoRoutinePool struct {
	mu       sync.RWMutex
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
	maxSize  int
	current  int
	taskChan chan func()
	logger   *zap.Logger
}

// NewGoRoutinePool 创建 goroutine 池
func NewGoRoutinePool(maxSize int, logger *zap.Logger) *GoRoutinePool {
	ctx, cancel := context.WithCancel(context.Background())
	pool := &GoRoutinePool{
		maxSize:  maxSize,
		taskChan: make(chan func(), maxSize*2), // 缓冲队列
		ctx:      ctx,
		cancel:   cancel,
		logger:   logger,
	}

	// 启动工作 goroutine
	for i := 0; i < maxSize; i++ {
		pool.wg.Add(1)
		go pool.worker()
	}

	return pool
}

// worker 工作 goroutine
func (p *GoRoutinePool) worker() {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return
		case task, ok := <-p.taskChan:
			if !ok {
				return
			}

			// 执行任务，捕获 panic
			func() {
				defer func() {
					if r := recover(); r != nil {
						if p.logger != nil {
							p.logger.Error("goroutine 任务发生 panic",
								zap.Any("panic", r),
								zap.String("stack", string(debug.Stack())),
							)
						}
					}
				}()
				task()
			}()
		}
	}
}

// Submit 提交任务到池
func (p *GoRoutinePool) Submit(task func()) error {
	select {
	case <-p.ctx.Done():
		return fmt.Errorf("goroutine pool is closed")
	case p.taskChan <- task:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("submit task timeout")
	}
}

// Shutdown 关闭 goroutine 池
func (p *GoRoutinePool) Shutdown(timeout time.Duration) error {
	close(p.taskChan)
	p.cancel()

	// 等待所有 goroutine 完成，带超时
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("shutdown timeout after %v", timeout)
	}
}

// Size 返回当前活跃的 goroutine 数量
func (p *GoRoutinePool) Size() int {
	return p.maxSize
}

// TransactionGuard 事务守卫，确保事务正确提交或回滚
type TransactionGuard struct {
	tx       interface{ Rollback() error }
	logger   *zap.Logger
	rolledBack bool
	committed  bool
}

// NewTransactionGuard 创建事务守卫
func NewTransactionGuard(tx interface{ Rollback() error }, logger *zap.Logger) *TransactionGuard {
	return &TransactionGuard{
		tx:     tx,
		logger: logger,
	}
}

// Commit 提交事务
func (g *TransactionGuard) Commit() {
	g.committed = true
}

// Rollback 回滚事务
func (g *TransactionGuard) Rollback() error {
	if g.rolledBack || g.committed {
		return nil
	}
	g.rolledBack = true
	return g.tx.Rollback()
}

// Release 释放事务（defer 调用）
func (g *TransactionGuard) Release() {
	if g.committed || g.rolledBack {
		return
	}

	// 既没有提交也没有回滚，自动回滚
	if err := g.tx.Rollback(); err != nil {
		if g.logger != nil {
			g.logger.Error("事务自动回滚失败",
				zap.Error(err),
				zap.String("stack", string(debug.Stack())),
			)
		}
	} else {
		if g.logger != nil {
			g.logger.Warn("事务未正确提交，已自动回滚",
				zap.String("stack", string(debug.Stack())),
			)
		}
	}
}

// ContextGuard Context 守卫，确保 cancel 被调用
type ContextGuard struct {
	ctx       context.Context
	cancel    context.CancelFunc
	logger    *zap.Logger
	cancelled bool
}

// NewContextGuard 创建 Context 守卫
func NewContextGuard(ctx context.Context, cancel context.CancelFunc, logger *zap.Logger) *ContextGuard {
	return &ContextGuard{
		ctx:    ctx,
		cancel: cancel,
		logger: logger,
	}
}

// Cancel 取消 context
func (g *ContextGuard) Cancel() {
	if !g.cancelled {
		g.cancelled = true
		g.cancel()
	}
}

// Release 释放 context（defer 调用）
func (g *ContextGuard) Release() {
	if !g.cancelled {
		g.cancel()
		if g.logger != nil {
			g.logger.Debug("Context 自动取消",
				zap.String("stack", string(debug.Stack())),
			)
		}
	}
}
