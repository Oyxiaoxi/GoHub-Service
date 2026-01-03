package singleflight

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	var g Group
	v, err := g.Do("key", func() (interface{}, error) {
		return "bar", nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "bar", v)
}

func TestDoErr(t *testing.T) {
	var g Group
	someErr := errors.New("some error")
	v, err := g.Do("key", func() (interface{}, error) {
		return nil, someErr
	})

	assert.Equal(t, someErr, err)
	assert.Nil(t, v)
}

func TestDoDupSuppress(t *testing.T) {
	var g Group
	var calls int32

	// 启动10个并发调用
	var wg sync.WaitGroup
	const n = 10
	wg.Add(n)
	
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			v, err := g.Do("key", func() (interface{}, error) {
				atomic.AddInt32(&calls, 1)
				time.Sleep(10 * time.Millisecond) // 模拟耗时操作
				return "bar", nil
			})
			assert.NoError(t, err)
			assert.Equal(t, "bar", v)
		}()
	}

	wg.Wait()
	
	// 虽然有10个并发调用，但函数只应该执行1次
	assert.Equal(t, int32(1), atomic.LoadInt32(&calls))
}

func TestForget(t *testing.T) {
	var g Group

	var firstDone = make(chan struct{})
	var calls int32

	// 第一次调用
	go func() {
		g.Do("key", func() (interface{}, error) {
			atomic.AddInt32(&calls, 1)
			<-firstDone
			return "bar", nil
		})
	}()

	// 等待第一次调用开始
	time.Sleep(10 * time.Millisecond)

	// Forget这个key
	g.Forget("key")

	// 第二次调用应该会再次执行函数
	v, err := g.Do("key", func() (interface{}, error) {
		atomic.AddInt32(&calls, 1)
		return "baz", nil
	})

	close(firstDone)

	assert.NoError(t, err)
	assert.Equal(t, "baz", v)
	
	// 等待确保第一次调用完成
	time.Sleep(50 * time.Millisecond)
	
	// 应该有2次调用（因为用了Forget）
	assert.Equal(t, int32(2), atomic.LoadInt32(&calls))
}

func BenchmarkDo(b *testing.B) {
	var g Group
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			g.Do("key", func() (interface{}, error) {
				return "bar", nil
			})
		}
	})
}
