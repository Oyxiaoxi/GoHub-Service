// Package cache 缓存降级策略
package cache

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var (
	// ErrCacheDegraded 缓存已降级
	ErrCacheDegraded = errors.New("cache degraded")
)

// DegradationConfig 降级配置
type DegradationConfig struct {
	// 是否启用自动降级
	EnableAutoDegradation bool
	
	// 错误率阈值（%）超过此值触发降级
	ErrorRateThreshold float64
	
	// 统计窗口时间
	WindowDuration time.Duration
	
	// 降级持续时间
	DegradationDuration time.Duration
	
	// 最小请求数（达到此数量才开始统计错误率）
	MinRequestCount int64
}

// DefaultDegradationConfig 默认降级配置
func DefaultDegradationConfig() *DegradationConfig {
	return &DegradationConfig{
		EnableAutoDegradation: true,
		ErrorRateThreshold:    50.0, // 错误率超过50%触发降级
		WindowDuration:        1 * time.Minute,
		DegradationDuration:   5 * time.Minute,
		MinRequestCount:       10,
	}
}

// DegradationManager 降级管理器
type DegradationManager struct {
	config *DegradationConfig
	
	// 降级状态
	degraded atomic.Bool
	degradedAt time.Time
	degradedMu sync.RWMutex
	
	// 统计数据
	totalRequests atomic.Int64
	errorRequests atomic.Int64
	windowStart   time.Time
	windowMu      sync.RWMutex
}

// NewDegradationManager 创建降级管理器
func NewDegradationManager(config *DegradationConfig) *DegradationManager {
	if config == nil {
		config = DefaultDegradationConfig()
	}
	
	dm := &DegradationManager{
		config:      config,
		windowStart: time.Now(),
	}
	
	// 启动后台任务
	if config.EnableAutoDegradation {
		go dm.autoCheck()
	}
	
	return dm
}

// RecordRequest 记录请求（成功）
func (dm *DegradationManager) RecordRequest() {
	dm.totalRequests.Add(1)
}

// RecordError 记录错误请求
func (dm *DegradationManager) RecordError() {
	dm.totalRequests.Add(1)
	dm.errorRequests.Add(1)
}

// IsDegraded 检查是否处于降级状态
func (dm *DegradationManager) IsDegraded() bool {
	if !dm.config.EnableAutoDegradation {
		return false
	}
	
	// 检查降级状态
	if dm.degraded.Load() {
		dm.degradedMu.RLock()
		degradedAt := dm.degradedAt
		dm.degradedMu.RUnlock()
		
		// 检查是否应该恢复
		if time.Since(degradedAt) > dm.config.DegradationDuration {
			dm.Recover()
			return false
		}
		return true
	}
	
	return false
}

// Degrade 手动触发降级
func (dm *DegradationManager) Degrade() {
	dm.degraded.Store(true)
	dm.degradedMu.Lock()
	dm.degradedAt = time.Now()
	dm.degradedMu.Unlock()
}

// Recover 恢复正常
func (dm *DegradationManager) Recover() {
	dm.degraded.Store(false)
	dm.resetStats()
}

// GetErrorRate 获取当前错误率
func (dm *DegradationManager) GetErrorRate() float64 {
	total := dm.totalRequests.Load()
	if total < dm.config.MinRequestCount {
		return 0
	}
	
	errors := dm.errorRequests.Load()
	return float64(errors) / float64(total) * 100
}

// GetStats 获取统计信息
func (dm *DegradationManager) GetStats() map[string]interface{} {
	dm.windowMu.RLock()
	windowStart := dm.windowStart
	dm.windowMu.RUnlock()
	
	dm.degradedMu.RLock()
	degradedAt := dm.degradedAt
	dm.degradedMu.RUnlock()
	
	return map[string]interface{}{
		"degraded":        dm.IsDegraded(),
		"error_rate":      dm.GetErrorRate(),
		"total_requests":  dm.totalRequests.Load(),
		"error_requests":  dm.errorRequests.Load(),
		"window_start":    windowStart,
		"degraded_at":     degradedAt,
		"window_duration": dm.config.WindowDuration,
	}
}

// autoCheck 自动检查是否需要降级
func (dm *DegradationManager) autoCheck() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	
	for range ticker.C {
		dm.checkAndDegrade()
		dm.rotateWindow()
	}
}

// checkAndDegrade 检查并触发降级
func (dm *DegradationManager) checkAndDegrade() {
	// 如果已经降级，跳过检查
	if dm.degraded.Load() {
		return
	}
	
	// 检查错误率
	errorRate := dm.GetErrorRate()
	if errorRate >= dm.config.ErrorRateThreshold {
		dm.Degrade()
	}
}

// rotateWindow 轮转统计窗口
func (dm *DegradationManager) rotateWindow() {
	dm.windowMu.RLock()
	windowStart := dm.windowStart
	dm.windowMu.RUnlock()
	
	// 检查是否需要重置窗口
	if time.Since(windowStart) >= dm.config.WindowDuration {
		dm.resetStats()
	}
}

// resetStats 重置统计数据
func (dm *DegradationManager) resetStats() {
	dm.totalRequests.Store(0)
	dm.errorRequests.Store(0)
	dm.windowMu.Lock()
	dm.windowStart = time.Now()
	dm.windowMu.Unlock()
}

// 全局降级管理器
var (
	globalDegradation     *DegradationManager
	globalDegradationOnce sync.Once
)

// GetDegradationManager 获取全局降级管理器
func GetDegradationManager() *DegradationManager {
	globalDegradationOnce.Do(func() {
		globalDegradation = NewDegradationManager(DefaultDegradationConfig())
	})
	return globalDegradation
}

// SafeGetWithDegradation 带降级的安全获取
func SafeGetWithDegradation(
	ctx context.Context,
	key string,
	ttl time.Duration,
	fetchFunc func() (interface{}, error),
	fallbackFunc func() (interface{}, error),
) (interface{}, error) {
	dm := GetDegradationManager()
	
	// 检查是否降级
	if dm.IsDegraded() {
		dm.RecordRequest()
		// 降级状态：直接返回降级数据
		if fallbackFunc != nil {
			return fallbackFunc()
		}
		return nil, ErrCacheDegraded
	}
	
	// 尝试从缓存获取
	value, err := SafeGet(ctx, key, ttl, fetchFunc)
	
	// 记录统计
	if err != nil {
		dm.RecordError()
		// 发生错误时，尝试使用降级函数
		if fallbackFunc != nil {
			return fallbackFunc()
		}
		return nil, err
	}
	
	dm.RecordRequest()
	return value, nil
}

// WarmupTask 预热任务
type WarmupTask struct {
	Name        string
	Description string
	Priority    int // 优先级，数字越小优先级越高
	Enabled     bool
	Fn          func(context.Context) error
}

// WarmupScheduler 预热调度器
type WarmupScheduler struct {
	tasks []*WarmupTask
	mu    sync.RWMutex
}

var (
	globalScheduler     *WarmupScheduler
	globalSchedulerOnce sync.Once
)

// GetScheduler 获取全局调度器
func GetScheduler() *WarmupScheduler {
	globalSchedulerOnce.Do(func() {
		globalScheduler = &WarmupScheduler{
			tasks: make([]*WarmupTask, 0),
		}
	})
	return globalScheduler
}

// RegisterTask 注册预热任务
func (ws *WarmupScheduler) RegisterTask(task *WarmupTask) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.tasks = append(ws.tasks, task)
}

// ExecuteAll 执行所有预热任务
func (ws *WarmupScheduler) ExecuteAll(ctx context.Context) []error {
	ws.mu.RLock()
	tasks := make([]*WarmupTask, len(ws.tasks))
	copy(tasks, ws.tasks)
	ws.mu.RUnlock()
	
	// 按优先级排序
	for i := 0; i < len(tasks); i++ {
		for j := i + 1; j < len(tasks); j++ {
			if tasks[i].Priority > tasks[j].Priority {
				tasks[i], tasks[j] = tasks[j], tasks[i]
			}
		}
	}
	
	// 执行任务
	var errors []error
	for _, task := range tasks {
		if !task.Enabled {
			continue
		}
		
		if err := task.Fn(ctx); err != nil {
			errors = append(errors, err)
		}
	}
	
	return errors
}

// ExecuteOne 执行指定任务
func (ws *WarmupScheduler) ExecuteOne(ctx context.Context, name string) error {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	
	for _, task := range ws.tasks {
		if task.Name == name {
			if !task.Enabled {
				return errors.New("task is disabled")
			}
			return task.Fn(ctx)
		}
	}
	
	return errors.New("task not found")
}

// ListTasks 列出所有任务
func (ws *WarmupScheduler) ListTasks() []*WarmupTask {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	
	tasks := make([]*WarmupTask, len(ws.tasks))
	copy(tasks, ws.tasks)
	return tasks
}

// RegisterWarmupTask 注册预热任务（便捷方法）
func RegisterWarmupTask(name, description string, priority int, fn func(context.Context) error) {
	GetScheduler().RegisterTask(&WarmupTask{
		Name:        name,
		Description: description,
		Priority:    priority,
		Enabled:     true,
		Fn:          fn,
	})
}

// WarmupAllTasks 执行所有预热任务（便捷方法）
func WarmupAllTasks(ctx context.Context) []error {
	return GetScheduler().ExecuteAll(ctx)
}
