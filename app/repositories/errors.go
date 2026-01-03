// Package repositories 仓储层错误定义
package repositories

import (
	apperrors "GoHub-Service/pkg/errors"
)

// 仓储层预定义错误（用于errors.Is比较）
var (
	// ErrNotFound 资源不存在
	ErrNotFound = apperrors.NotFoundError("资源")
	
	// ErrCreateFailed 创建失败
	ErrCreateFailed = apperrors.DatabaseCreateError("资源", nil)
	
	// ErrUpdateFailed 更新失败
	ErrUpdateFailed = apperrors.DatabaseUpdateError("资源", nil)
	
	// ErrDeleteFailed 删除失败
	ErrDeleteFailed = apperrors.DatabaseDeleteError("资源", nil)
	
	// ErrDuplicateKey 重复键错误
	ErrDuplicateKey = apperrors.DatabaseDuplicateError("资源")
	
	// ErrQueryFailed 查询失败
	ErrQueryFailed = apperrors.DatabaseError("查询", nil)
)

// NewNotFoundError 创建资源不存在错误
func NewNotFoundError(resource string, id interface{}) *apperrors.AppError {
	err := apperrors.NotFoundError(resource)
	err.WithDetails(map[string]interface{}{
		"resource": resource,
		"id":       id,
	})
	return err
}

// NewCreateError 创建失败错误
func NewCreateError(resource string, err error) *apperrors.AppError {
	return apperrors.DatabaseCreateError(resource, err)
}

// NewUpdateError 更新失败错误
func NewUpdateError(resource string, id interface{}, err error) *apperrors.AppError {
	e := apperrors.DatabaseUpdateError(resource, err)
	e.WithDetails(map[string]interface{}{
		"resource": resource,
		"id":       id,
	})
	return e
}

// NewDeleteError 删除失败错误
func NewDeleteError(resource string, id interface{}, err error) *apperrors.AppError {
	e := apperrors.DatabaseDeleteError(resource, err)
	e.WithDetails(map[string]interface{}{
		"resource": resource,
		"id":       id,
	})
	return e
}

// NewQueryError 查询失败错误
func NewQueryError(operation string, err error) *apperrors.AppError {
	return apperrors.DatabaseError(operation, err)
}

// NewDuplicateError 重复记录错误
func NewDuplicateError(resource string, field string, value interface{}) *apperrors.AppError {
	err := apperrors.DatabaseDuplicateError(resource)
	err.WithDetails(map[string]interface{}{
		"resource": resource,
		"field":    field,
		"value":    value,
	})
	return err
}
