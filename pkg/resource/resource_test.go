package resource

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// 模拟可关闭的资源
type mockCloser struct {
	closed bool
	err    error
}

func (m *mockCloser) Close() error {
	m.closed = true
	return m.err
}

func TestSafeClose(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("正常关闭", func(t *testing.T) {
		closer := &mockCloser{}
		SafeClose(closer, logger)
		assert.True(t, closer.closed)
	})

	t.Run("关闭失败", func(t *testing.T) {
		closer := &mockCloser{err: errors.New("close error")}
		SafeClose(closer, logger)
		assert.True(t, closer.closed)
	})

	t.Run("nil closer", func(t *testing.T) {
		SafeClose(nil, logger)
		// 不应该 panic
	})
}

func TestTracker(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	tracker := NewTracker(logger)

	t.Run("追踪和取消追踪", func(t *testing.T) {
		tracker.Track("res1", "http.Response")
		assert.Equal(t, 1, tracker.Count())

		tracker.Untrack("res1")
		assert.Equal(t, 0, tracker.Count())
	})

	t.Run("检测泄漏", func(t *testing.T) {
		tracker.Clear()

		// 添加一个旧资源
		tracker.Track("res2", "connection")
		time.Sleep(100 * time.Millisecond)

		leaked := tracker.Check(50 * time.Millisecond)
		assert.Len(t, leaked, 1)
		assert.Equal(t, "res2", leaked[0].ID)
	})

	t.Run("报告泄漏", func(t *testing.T) {
		tracker.Clear()
		tracker.Track("res3", "file")
		time.Sleep(100 * time.Millisecond)

		tracker.Report(50 * time.Millisecond)
		// 应该记录警告日志
	})
}

func TestGoRoutinePool(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("正常提交任务", func(t *testing.T) {
		pool := NewGoRoutinePool(2, logger)
		defer pool.Shutdown(time.Second)

		executed := false
		err := pool.Submit(func() {
			executed = true
		})

		assert.NoError(t, err)
		time.Sleep(100 * time.Millisecond)
		assert.True(t, executed)
	})

	t.Run("任务 panic 不影响池", func(t *testing.T) {
		pool := NewGoRoutinePool(2, logger)
		defer pool.Shutdown(time.Second)

		err := pool.Submit(func() {
			panic("test panic")
		})

		assert.NoError(t, err)
		time.Sleep(100 * time.Millisecond)

		// 池应该仍然可用
		executed := false
		err = pool.Submit(func() {
			executed = true
		})

		assert.NoError(t, err)
		time.Sleep(100 * time.Millisecond)
		assert.True(t, executed)
	})

	t.Run("关闭后不能提交", func(t *testing.T) {
		pool := NewGoRoutinePool(2, logger)
		pool.Shutdown(time.Second)

		err := pool.Submit(func() {})
		assert.Error(t, err)
	})

	t.Run("正常关闭", func(t *testing.T) {
		pool := NewGoRoutinePool(2, logger)

		// 提交一些任务
		for i := 0; i < 5; i++ {
			pool.Submit(func() {
				time.Sleep(50 * time.Millisecond)
			})
		}

		err := pool.Shutdown(2 * time.Second)
		assert.NoError(t, err)
	})
}

func TestTransactionGuard(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	type mockTx struct {
		rolledBack bool
	}

	t.Run("正常提交", func(t *testing.T) {
		tx := &mockTx{}
		guard := NewTransactionGuard(&struct{ Rollback func() error }{
			Rollback: func() error {
				tx.rolledBack = true
				return nil
			},
		}, logger)

		guard.Commit()
		guard.Release()

		// 不应该回滚
		assert.False(t, tx.rolledBack)
	})

	t.Run("手动回滚", func(t *testing.T) {
		tx := &mockTx{}
		guard := NewTransactionGuard(&struct{ Rollback func() error }{
			Rollback: func() error {
				tx.rolledBack = true
				return nil
			},
		}, logger)

		guard.Rollback()
		guard.Release()

		assert.True(t, tx.rolledBack)
	})

	t.Run("自动回滚", func(t *testing.T) {
		tx := &mockTx{}
		guard := NewTransactionGuard(&struct{ Rollback func() error }{
			Rollback: func() error {
				tx.rolledBack = true
				return nil
			},
		}, logger)

		// 既没有提交也没有回滚
		guard.Release()

		// 应该自动回滚
		assert.True(t, tx.rolledBack)
	})
}

func TestContextGuard(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("手动取消", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		guard := NewContextGuard(ctx, cancel, logger)

		guard.Cancel()
		guard.Release()

		select {
		case <-ctx.Done():
			// 正常
		default:
			t.Fatal("context should be cancelled")
		}
	})

	t.Run("自动取消", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		guard := NewContextGuard(ctx, cancel, logger)

		// 不手动取消
		guard.Release()

		select {
		case <-ctx.Done():
			// 正常，应该自动取消
		default:
			t.Fatal("context should be cancelled automatically")
		}
	})
}

// Benchmark 测试
func BenchmarkSafeClose(b *testing.B) {
	logger, _ := zap.NewDevelopment()
	closer := &mockCloser{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SafeClose(closer, logger)
	}
}

func BenchmarkTracker(b *testing.B) {
	logger, _ := zap.NewDevelopment()
	tracker := NewTracker(logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := time.Now().String()
		tracker.Track(id, "test")
		tracker.Untrack(id)
	}
}

func BenchmarkGoRoutinePool(b *testing.B) {
	logger, _ := zap.NewDevelopment()
	pool := NewGoRoutinePool(10, logger)
	defer pool.Shutdown(time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pool.Submit(func() {
			time.Sleep(1 * time.Millisecond)
		})
	}
}
