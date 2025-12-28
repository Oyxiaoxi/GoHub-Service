// Package repositories 仓储层错误定义
package repositories

import "errors"

var (
	// ErrNotFound 资源不存在
	ErrNotFound = errors.New("resource not found")
	
	// ErrCreateFailed 创建失败
	ErrCreateFailed = errors.New("create failed")
	
	// ErrUpdateFailed 更新失败
	ErrUpdateFailed = errors.New("update failed")
	
	// ErrDeleteFailed 删除失败
	ErrDeleteFailed = errors.New("delete failed")
)
