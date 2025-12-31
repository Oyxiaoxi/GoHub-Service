package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// SlowQuery 慢查询记录结构
type SlowQuery struct {
	Query        string
	ExecutionTime float64
	LockTime     float64
	RowsSent     int
	RowsExamined int
	Timestamp    time.Time
}

// SlowQueryStats 慢查询统计
type SlowQueryStats struct {
	Query           string
	Count           int
	TotalTime       float64
	AvgTime         float64
	MaxTime         float64
	MinTime         float64
	TotalRowsExam   int
	AvgRowsExam     float64
}

var CmdSlowLog = &cobra.Command{
	Use:   "slowlog",
	Short: "Analyze MySQL slow query log",
	Run:   runSlowLog,
}

func init() {
	CmdSlowLog.Flags().StringP("file", "f", "", "Slow query log file path")
	CmdSlowLog.Flags().IntP("top", "t", 10, "Show top N slow queries")
	CmdSlowLog.Flags().Float64P("threshold", "s", 0.0, "Filter queries slower than threshold (seconds)")
}

func runSlowLog(cmd *cobra.Command, args []string) {
	logFile, _ := cmd.Flags().GetString("file")
	topN, _ := cmd.Flags().GetInt("top")
	threshold, _ := cmd.Flags().GetFloat64("threshold")

	if logFile == "" {
		fmt.Println("❌ Please specify slow query log file with -f flag")
		fmt.Println("Example: ./main slowlog -f /var/log/mysql/slow-query.log")
		return
	}

	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		fmt.Printf("❌ Log file not found: %s\n", logFile)
		return
	}

	fmt.Printf("📊 Analyzing slow query log: %s\n", logFile)
	fmt.Println(strings.Repeat("=", 80))

	queries, err := parseSlowQueryLog(logFile, threshold)
	if err != nil {
		fmt.Printf("❌ Failed to parse log: %v\n", err)
		return
	}

	if len(queries) == 0 {
		fmt.Println("✅ No slow queries found!")
		return
	}

	stats := analyzeQueries(queries)
	printSummary(queries, threshold)
	printTopQueries(stats, topN)
	
	fmt.Println("\n✅ Analysis completed!")
}

func parseSlowQueryLog(filePath string, threshold float64) ([]SlowQuery, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var queries []SlowQuery
	var currentQuery SlowQuery
	var queryBuffer strings.Builder
	inQuery := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "# Query_time:") {
			parts := strings.Fields(line)
			for i, part := range parts {
				if part == "Query_time:" && i+1 < len(parts) {
					fmt.Sscanf(parts[i+1], "%f", &currentQuery.ExecutionTime)
				} else if part == "Rows_examined:" && i+1 < len(parts) {
					fmt.Sscanf(parts[i+1], "%d", &currentQuery.RowsExamined)
				}
			}
		} else if strings.HasPrefix(strings.ToUpper(line), "SELECT") ||
			strings.HasPrefix(strings.ToUpper(line), "UPDATE") ||
			strings.HasPrefix(strings.ToUpper(line), "INSERT") ||
			strings.HasPrefix(strings.ToUpper(line), "DELETE") {
			inQuery = true
			queryBuffer.WriteString(line)
			queryBuffer.WriteString(" ")
		} else if inQuery {
			if strings.HasSuffix(line, ";") {
				queryBuffer.WriteString(line)
				currentQuery.Query = queryBuffer.String()

				if threshold == 0 || currentQuery.ExecutionTime >= threshold {
					queries = append(queries, currentQuery)
				}

				queryBuffer.Reset()
				currentQuery = SlowQuery{}
				inQuery = false
			} else if line != "" {
				queryBuffer.WriteString(line)
				queryBuffer.WriteString(" ")
			}
		}
	}

	return queries, nil
}

func analyzeQueries(queries []SlowQuery) []SlowQueryStats {
	statsMap := make(map[string]*SlowQueryStats)

	for _, q := range queries {
		normalizedQuery := normalizeQuery(q.Query)

		if stat, exists := statsMap[normalizedQuery]; exists {
			stat.Count++
			stat.TotalTime += q.ExecutionTime
			stat.TotalRowsExam += q.RowsExamined

			if q.ExecutionTime > stat.MaxTime {
				stat.MaxTime = q.ExecutionTime
			}
			if q.ExecutionTime < stat.MinTime {
				stat.MinTime = q.ExecutionTime
			}
		} else {
			statsMap[normalizedQuery] = &SlowQueryStats{
				Query:         normalizedQuery,
				Count:         1,
				TotalTime:     q.ExecutionTime,
				MaxTime:       q.ExecutionTime,
				MinTime:       q.ExecutionTime,
				TotalRowsExam: q.RowsExamined,
			}
		}
	}

	var stats []SlowQueryStats
	for _, stat := range statsMap {
		stat.AvgTime = stat.TotalTime / float64(stat.Count)
		stat.AvgRowsExam = float64(stat.TotalRowsExam) / float64(stat.Count)
		stats = append(stats, *stat)
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].TotalTime > stats[j].TotalTime
	})

	return stats
}

func normalizeQuery(query string) string {
	query = strings.TrimSpace(query)
	query = strings.TrimSuffix(query, ";")
	// 简单标准化
	return query
}

func printSummary(queries []SlowQuery, threshold float64) {
	totalTime := 0.0
	totalRowsExam := 0

	for _, q := range queries {
		totalTime += q.ExecutionTime
		totalRowsExam += q.RowsExamined
	}

	fmt.Println("\n📊 Summary Statistics")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("Total slow queries: %d\n", len(queries))
	if threshold > 0 {
		fmt.Printf("Threshold: >= %.2f seconds\n", threshold)
	}
	fmt.Printf("Total execution time: %.2f seconds\n", totalTime)
	fmt.Printf("Average execution time: %.4f seconds\n", totalTime/float64(len(queries)))
	fmt.Printf("Total rows examined: %d\n", totalRowsExam)
	fmt.Printf("Average rows examined: %.0f\n", float64(totalRowsExam)/float64(len(queries)))
}

func printTopQueries(stats []SlowQueryStats, topN int) {
	fmt.Printf("\n🔍 Top %d Slow Queries\n", topN)
	fmt.Println(strings.Repeat("=", 80))

	if topN > len(stats) {
		topN = len(stats)
	}

	for i := 0; i < topN; i++ {
		stat := stats[i]
		fmt.Printf("\n#%d\n", i+1)
		fmt.Println(strings.Repeat("-", 80))
		fmt.Printf("Count: %d times\n", stat.Count)
		fmt.Printf("Total Time: %.4fs | Avg: %.4fs | Max: %.4fs | Min: %.4fs\n",
			stat.TotalTime, stat.AvgTime, stat.MaxTime, stat.MinTime)
		fmt.Printf("Avg Rows Examined: %.0f\n", stat.AvgRowsExam)
		
		query := stat.Query
		if len(query) > 150 {
			query = query[:150] + "..."
		}
		fmt.Printf("\nQuery: %s\n", query)

		// 优化建议
		if stat.AvgRowsExam > 10000 {
			fmt.Println("💡 Suggestion: High rows examined (>10K). Consider adding indexes.")
		} else if stat.AvgRowsExam > 1000 {
			fmt.Println("💡 Suggestion: Moderate rows examined. Check index usage.")
		}
		if stat.AvgTime > 1.0 {
			fmt.Println("⚠️  Very slow query (>1s). Urgent optimization needed.")
		}
		if stat.Count > 100 {
			fmt.Println("💡 Suggestion: High frequency query. Consider caching.")
		}
	}
}
