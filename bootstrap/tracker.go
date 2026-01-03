package bootstrap

import (
	"GoHub-Service/pkg/logger"
	"GoHub-Service/pkg/resource"
	"time"

	"go.uber.org/zap"
)

// Tracker 全局资源追踪器
var Tracker *resource.Tracker

// SetupTracker 初始化资源追踪器
func SetupTracker() {
	Tracker = resource.NewTracker(logger.Logger)
}

// StartTrackerReporting 启动资源泄漏定期报告
// threshold: 资源超过该时间未释放会被认为可能泄漏（建议5分钟）
// interval: 检查间隔（建议1分钟）
func StartTrackerReporting(threshold, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			Tracker.Report(threshold)

			// 记录当前追踪的资源数量
			count := Tracker.Count()
			if count > 0 {
				logger.Logger.Info("资源追踪统计",
					zap.Int("tracked_count", count),
				)
			}
		}
	}()
}
