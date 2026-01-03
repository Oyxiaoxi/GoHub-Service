// Package database 慢查询日志记录
package database

import (
	"context"
	"time"

	"GoHub-Service/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// SlowQueryThreshold 慢查询阈值（毫秒）
const SlowQueryThreshold = 200 * time.Millisecond

// SlowQueryLogger 慢查询日志记录器
type SlowQueryLogger struct {
	gormlogger.Interface
	threshold time.Duration
}

// NewSlowQueryLogger 创建慢查询日志记录器
func NewSlowQueryLogger(threshold time.Duration) *SlowQueryLogger {
	if threshold == 0 {
		threshold = SlowQueryThreshold
	}
	return &SlowQueryLogger{
		Interface: gormlogger.Default.LogMode(gormlogger.Info),
		threshold: threshold,
	}
}

// Trace 记录 SQL 执行时间
func (l *SlowQueryLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	
	// 只记录慢查询
	if elapsed >= l.threshold {
		sql, rows := fc()
		
		logger.WarnContext(ctx, "slow_query",
			zap.String("sql", sql),
			zap.Int64("elapsed_ms", elapsed.Milliseconds()),
			zap.Int64("rows_affected", rows),
			zap.Int64("threshold_ms", l.threshold.Milliseconds()),
		)
	}
	
	// 调用原始 Trace 方法
	if l.Interface != nil {
		l.Interface.Trace(ctx, begin, fc, err)
	}
}

// EnableSlowQueryLog 启用慢查询日志
func EnableSlowQueryLog(db *gorm.DB, threshold time.Duration) *gorm.DB {
	slowLogger := NewSlowQueryLogger(threshold)
	return db.Session(&gorm.Session{
		Logger: slowLogger,
	})
}

// SlowQueryStats 慢查询统计
type SlowQueryStats struct {
	TotalCount    int64         `json:"total_count"`
	SlowCount     int64         `json:"slow_count"`
	AverageTime   time.Duration `json:"average_time"`
	MaxTime       time.Duration `json:"max_time"`
	SlowQueries   []SlowQuery   `json:"slow_queries"`
}

// SlowQuery 慢查询记录
type SlowQuery struct {
	SQL       string        `json:"sql"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
	Rows      int64         `json:"rows"`
}

// slowQueryRecorder 全局慢查询记录器
var slowQueryRecorder = &struct {
	queries []SlowQuery
	maxSize int
}{
	queries: make([]SlowQuery, 0, 100),
	maxSize: 100,
}

// RecordSlowQuery 记录慢查询
func RecordSlowQuery(sql string, duration time.Duration, rows int64) {
	if len(slowQueryRecorder.queries) >= slowQueryRecorder.maxSize {
		// 移除最旧的记录
		slowQueryRecorder.queries = slowQueryRecorder.queries[1:]
	}
	
	slowQueryRecorder.queries = append(slowQueryRecorder.queries, SlowQuery{
		SQL:       sql,
		Duration:  duration,
		Timestamp: time.Now(),
		Rows:      rows,
	})
}

// GetSlowQueryStats 获取慢查询统计
func GetSlowQueryStats() SlowQueryStats {
	stats := SlowQueryStats{
		SlowCount:   int64(len(slowQueryRecorder.queries)),
		SlowQueries: slowQueryRecorder.queries,
	}
	
	if stats.SlowCount > 0 {
		var total time.Duration
		maxTime := time.Duration(0)
		
		for _, q := range slowQueryRecorder.queries {
			total += q.Duration
			if q.Duration > maxTime {
				maxTime = q.Duration
			}
		}
		
		stats.AverageTime = total / time.Duration(stats.SlowCount)
		stats.MaxTime = maxTime
	}
	
	return stats
}

// ClearSlowQueryStats 清除慢查询统计
func ClearSlowQueryStats() {
	slowQueryRecorder.queries = make([]SlowQuery, 0, 100)
}
