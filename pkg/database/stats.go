// Package database 数据库连接池统计
package database

import (
	"fmt"
	"time"
)

// Stats 数据库连接池统计信息
type Stats struct {
	// MaxOpenConnections 最大打开连接数
	MaxOpenConnections int `json:"max_open_connections"`
	
	// OpenConnections 当前打开的连接数（使用中+空闲）
	OpenConnections int `json:"open_connections"`
	
	// InUse 正在使用的连接数
	InUse int `json:"in_use"`
	
	// Idle 空闲连接数
	Idle int `json:"idle"`
	
	// WaitCount 等待连接的总次数
	WaitCount int64 `json:"wait_count"`
	
	// WaitDuration 等待连接的总时间
	WaitDuration time.Duration `json:"wait_duration"`
	
	// MaxIdleClosed 因超过最大空闲数而关闭的连接数
	MaxIdleClosed int64 `json:"max_idle_closed"`
	
	// MaxLifetimeClosed 因超过最大生命周期而关闭的连接数
	MaxLifetimeClosed int64 `json:"max_lifetime_closed"`
	
	// MaxIdleTimeClosed 因超过最大空闲时间而关闭的连接数
	MaxIdleTimeClosed int64 `json:"max_idle_time_closed"`
	
	// 额外计算的指标
	// UtilizationRate 连接池使用率（使用中/最大连接数）
	UtilizationRate float64 `json:"utilization_rate"`
	
	// IdleRate 空闲连接率（空闲/打开连接）
	IdleRate float64 `json:"idle_rate"`
	
	// AvgWaitDuration 平均等待时间
	AvgWaitDuration time.Duration `json:"avg_wait_duration"`
}

// GetStats 获取数据库连接池统计信息
func GetStats() *Stats {
	if SQLDB == nil {
		return nil
	}
	
	dbStats := SQLDB.Stats()
	stats := &Stats{
		MaxOpenConnections:  dbStats.MaxOpenConnections,
		OpenConnections:     dbStats.OpenConnections,
		InUse:               dbStats.InUse,
		Idle:                dbStats.Idle,
		WaitCount:           dbStats.WaitCount,
		WaitDuration:        dbStats.WaitDuration,
		MaxIdleClosed:       dbStats.MaxIdleClosed,
		MaxLifetimeClosed:   dbStats.MaxLifetimeClosed,
		MaxIdleTimeClosed:   dbStats.MaxIdleTimeClosed,
	}
	
	// 计算使用率
	if stats.MaxOpenConnections > 0 {
		stats.UtilizationRate = float64(stats.InUse) / float64(stats.MaxOpenConnections) * 100
	}
	
	// 计算空闲率
	if stats.OpenConnections > 0 {
		stats.IdleRate = float64(stats.Idle) / float64(stats.OpenConnections) * 100
	}
	
	// 计算平均等待时间
	if stats.WaitCount > 0 {
		stats.AvgWaitDuration = stats.WaitDuration / time.Duration(stats.WaitCount)
	}
	
	return stats
}

// PrintStats 打印数据库连接池统计信息（用于调试）
func PrintStats() {
	stats := GetStats()
	if stats == nil {
		fmt.Println("数据库未初始化")
		return
	}
	
	fmt.Println("========== 数据库连接池统计 ==========")
	fmt.Printf("最大连接数:         %d\n", stats.MaxOpenConnections)
	fmt.Printf("当前打开连接:       %d\n", stats.OpenConnections)
	fmt.Printf("使用中连接:         %d\n", stats.InUse)
	fmt.Printf("空闲连接:           %d\n", stats.Idle)
	fmt.Printf("连接使用率:         %.2f%%\n", stats.UtilizationRate)
	fmt.Printf("空闲连接率:         %.2f%%\n", stats.IdleRate)
	fmt.Println("---------- 等待统计 ----------")
	fmt.Printf("等待次数:           %d\n", stats.WaitCount)
	fmt.Printf("总等待时间:         %v\n", stats.WaitDuration)
	fmt.Printf("平均等待时间:       %v\n", stats.AvgWaitDuration)
	fmt.Println("---------- 关闭统计 ----------")
	fmt.Printf("空闲超时关闭:       %d\n", stats.MaxIdleClosed)
	fmt.Printf("生命周期超时关闭:   %d\n", stats.MaxLifetimeClosed)
	fmt.Printf("空闲时间超时关闭:   %d\n", stats.MaxIdleTimeClosed)
	fmt.Println("=====================================")
}

// CheckHealth 检查数据库连接池健康状态
// 返回 true 表示健康，false 表示需要关注
func CheckHealth() (bool, []string) {
	stats := GetStats()
	if stats == nil {
		return false, []string{"数据库未初始化"}
	}
	
	warnings := []string{}
	healthy := true
	
	// 检查1：连接使用率过高（> 80%）
	if stats.UtilizationRate > 80 {
		healthy = false
		warnings = append(warnings, fmt.Sprintf("连接使用率过高: %.2f%% (建议 < 80%%)", stats.UtilizationRate))
	}
	
	// 检查2：等待次数过多
	if stats.WaitCount > 1000 {
		healthy = false
		warnings = append(warnings, fmt.Sprintf("等待次数过多: %d (建议增加最大连接数)", stats.WaitCount))
	}
	
	// 检查3：平均等待时间过长（> 100ms）
	if stats.AvgWaitDuration > 100*time.Millisecond {
		healthy = false
		warnings = append(warnings, fmt.Sprintf("平均等待时间过长: %v (建议 < 100ms)", stats.AvgWaitDuration))
	}
	
	// 检查4：空闲连接过多（> 50%）
	if stats.IdleRate > 50 && stats.OpenConnections > 10 {
		warnings = append(warnings, fmt.Sprintf("空闲连接过多: %.2f%% (可考虑降低最大空闲连接数)", stats.IdleRate))
	}
	
	// 检查5：连接频繁关闭（生命周期超时）
	if stats.MaxLifetimeClosed > 10000 {
		warnings = append(warnings, fmt.Sprintf("连接频繁关闭: %d (可考虑延长 max_life_seconds)", stats.MaxLifetimeClosed))
	}
	
	if len(warnings) == 0 {
		warnings = append(warnings, "连接池运行正常")
	}
	
	return healthy, warnings
}

// GetMetrics 获取 Prometheus 格式的指标（返回map便于集成）
func GetMetrics() map[string]interface{} {
	stats := GetStats()
	if stats == nil {
		return map[string]interface{}{
			"error": "database not initialized",
		}
	}
	
	return map[string]interface{}{
		"db_max_open_connections":      stats.MaxOpenConnections,
		"db_open_connections":          stats.OpenConnections,
		"db_in_use_connections":        stats.InUse,
		"db_idle_connections":          stats.Idle,
		"db_wait_count_total":          stats.WaitCount,
		"db_wait_duration_seconds":     stats.WaitDuration.Seconds(),
		"db_max_idle_closed_total":     stats.MaxIdleClosed,
		"db_max_lifetime_closed_total": stats.MaxLifetimeClosed,
		"db_max_idle_time_closed_total": stats.MaxIdleTimeClosed,
		"db_utilization_rate":          stats.UtilizationRate,
		"db_idle_rate":                 stats.IdleRate,
		"db_avg_wait_duration_seconds": stats.AvgWaitDuration.Seconds(),
	}
}

// RecommendConfig 根据当前统计推荐配置调整
func RecommendConfig() map[string]string {
	stats := GetStats()
	if stats == nil {
		return map[string]string{
			"error": "数据库未初始化",
		}
	}
	
	recommendations := make(map[string]string)
	
	// 推荐1：使用率过高，增加最大连接数
	if stats.UtilizationRate > 80 {
		newMax := stats.MaxOpenConnections * 2
		if newMax > 500 {
			newMax = 500
		}
		recommendations["max_open_connections"] = fmt.Sprintf("当前: %d, 建议: %d (使用率 %.2f%% 过高)", 
			stats.MaxOpenConnections, newMax, stats.UtilizationRate)
	}
	
	// 推荐2：空闲连接过多，减少最大空闲数
	if stats.IdleRate > 50 && stats.Idle > 20 {
		newIdle := stats.InUse + (stats.InUse / 2) // 使用中的1.5倍
		if newIdle < 10 {
			newIdle = 10
		}
		recommendations["max_idle_connections"] = fmt.Sprintf("当前空闲: %d, 建议: %d (空闲率 %.2f%% 过高)", 
			stats.Idle, newIdle, stats.IdleRate)
	}
	
	// 推荐3：连接频繁关闭，延长生命周期
	if stats.MaxLifetimeClosed > 10000 {
		recommendations["max_life_seconds"] = fmt.Sprintf("当前关闭数: %d, 建议延长至 20-30 分钟", 
			stats.MaxLifetimeClosed)
	}
	
	// 推荐4：等待时间过长
	if stats.AvgWaitDuration > 100*time.Millisecond {
		recommendations["performance"] = fmt.Sprintf("平均等待 %v，建议增加连接数或优化查询", 
			stats.AvgWaitDuration)
	}
	
	if len(recommendations) == 0 {
		recommendations["status"] = "连接池配置合理，无需调整"
	}
	
	return recommendations
}
