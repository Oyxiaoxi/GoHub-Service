// Package controller 控制器通用辅助方法
package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetIDFromParam 从URL参数中获取ID并转换为uint64
// 如果转换失败，返回0和false
func GetIDFromParam(c *gin.Context, paramName string) (uint64, bool) {
	idStr := c.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		return 0, false
	}
	return id, true
}

// GetIDParam 从URL参数中获取ID字符串
func GetIDParam(c *gin.Context) string {
	return c.Param("id")
}

// MustGetIDParam 获取ID参数，如果为空则返回false
func MustGetIDParam(c *gin.Context) (string, bool) {
	id := c.Param("id")
	if id == "" {
		return "", false
	}
	return id, true
}

// CheckModelID 检查模型ID是否有效（通过ID字段）
// 返回true表示有效，false表示无效
func CheckModelID(modelID interface{}) bool {
	// 检查ID是否为0或空字符串
	switch v := modelID.(type) {
	case uint, uint64, int, int64:
		return v != 0
	case string:
		return v != ""
	default:
		return false
	}
}

// CheckRowsAffected 检查数据库操作影响的行数
// 返回true表示操作成功，false表示失败
func CheckRowsAffected(rowsAffected int64) bool {
	return rowsAffected > 0
}
