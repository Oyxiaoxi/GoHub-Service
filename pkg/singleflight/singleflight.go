// Package singleflight 提供防止缓存击穿的单次执行机制
package singleflight

import (
	"sync"
)

// call 表示一个正在执行或已完成的函数调用
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// Group 管理一组函数调用，确保每个key只执行一次
type Group struct {
	mu sync.Mutex       // 保护m
	m  map[string]*call // 延迟初始化
}

// Do 执行给定的函数，确保对于相同的key只执行一次
// 如果有重复调用，重复的调用者会等待原始调用完成并接收相同的结果
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}

// Forget 告诉 singleflight 忘记一个key
// 未来对这个key的Do调用会调用函数而不是等待更早的调用完成
func (g *Group) Forget(key string) {
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
}
