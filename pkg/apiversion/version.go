// Package apiversion API 版本管理
package apiversion

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Version API 版本信息
type Version struct {
	Version     string   `json:"version"`
	Status      string   `json:"status"`       // active, deprecated, sunset
	ReleaseDate string   `json:"release_date"` // 发布日期
	SunsetDate  string   `json:"sunset_date"`  // 停用日期（可选）
	Features    []string `json:"features"`     // 主要特性
}

// APIVersions 所有支持的 API 版本
var APIVersions = map[string]Version{
	"v1": {
		Version:     "v1",
		Status:      "active",
		ReleaseDate: "2024-01-01",
		Features: []string{
			"用户管理",
			"话题管理",
			"评论管理",
			"分类管理",
			"权限管理",
			"搜索功能",
		},
	},
	"v2": {
		Version:     "v2",
		Status:      "planned",
		ReleaseDate: "2026-06-01",
		Features: []string{
			"GraphQL 支持",
			"Websocket 实时通知",
			"批量操作优化",
		},
	},
}

// CurrentVersion 当前默认版本
const CurrentVersion = "v1"

// GetVersionInfo 获取版本信息
func GetVersionInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"current_version": CurrentVersion,
			"versions":        APIVersions,
			"api_docs":        "/swagger/index.html",
		})
	}
}

// VersionDeprecated 版本废弃中间件
func VersionDeprecated(version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if v, exists := APIVersions[version]; exists {
			if v.Status == "deprecated" {
				c.Header("X-API-Warn", fmt.Sprintf("API version %s is deprecated. Sunset date: %s", version, v.SunsetDate))
			} else if v.Status == "sunset" {
				c.JSON(http.StatusGone, gin.H{
					"error": fmt.Sprintf("API version %s has been sunset on %s", version, v.SunsetDate),
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// RequireVersion 要求特定版本
func RequireVersion(minVersion string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestedVersion := c.Param("version")
		if requestedVersion == "" {
			requestedVersion = c.GetHeader("X-API-Version")
		}
		if requestedVersion == "" {
			requestedVersion = CurrentVersion
		}

		if !isValidVersion(requestedVersion) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Invalid API version: %s", requestedVersion),
			})
			c.Abort()
			return
		}

		c.Set("api_version", requestedVersion)
		c.Next()
	}
}

func isValidVersion(version string) bool {
	_, exists := APIVersions[version]
	return exists
}
