// Package testutil 提供测试数据工厂 - 简化版本
package testutil

import "time"

// MockTime 返回固定的测试时间
func MockTime() time.Time {
	return time.Date(2025, 12, 31, 12, 0, 0, 0, time.UTC)
}

// MockTimePtr 返回固定的测试时间指针
func MockTimePtr() *time.Time {
	t := MockTime()
	return &t
}
