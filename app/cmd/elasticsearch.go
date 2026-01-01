package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	"GoHub-Service/bootstrap"
	"GoHub-Service/pkg/elasticsearch"
)

// esCmd Elasticsearch 同步命令
var esCmd = &cobra.Command{
	Use:   "es",
	Short: "Elasticsearch related commands",
	Long:  "执行 Elasticsearch 数据同步和索引管理相关操作",
}

// esSyncCmd 完整同步命令
var esSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Full sync from MySQL to Elasticsearch",
	Long:  "完整同步MySQL数据库中的所有话题到Elasticsearch",
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化数据库连接
		bootstrap.LoadDB()
		bootstrap.LoadLogger()

		// 创建ES客户端
		addresses := []string{"http://localhost:9200"}
		client, err := elasticsearch.NewClient(addresses)
		if err != nil {
			log.Fatalf("Failed to create ES client: %v", err)
		}

		// 执行同步
		syncService := elasticsearch.NewSyncService(client)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		start := time.Now()
		if err := syncService.SyncAllTopics(ctx, 1000); err != nil {
			log.Fatalf("Sync failed: %v", err)
		}

		elapsed := time.Since(start)
		fmt.Printf("\n✓ Full sync completed in %v\n", elapsed)
	},
}

// esIncrementalCmd 增量同步命令
var esIncrementalCmd = &cobra.Command{
	Use:   "sync-incremental",
	Short: "Incremental sync from MySQL to Elasticsearch",
	Long:  "增量同步最近修改的话题到Elasticsearch",
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化数据库连接
		bootstrap.LoadDB()
		bootstrap.LoadLogger()

		// 创建ES客户端
		addresses := []string{"http://localhost:9200"}
		client, err := elasticsearch.NewClient(addresses)
		if err != nil {
			log.Fatalf("Failed to create ES client: %v", err)
		}

		// 获取时间间隔参数
		minutes, _ := cmd.Flags().GetInt("minutes")
		if minutes == 0 {
			minutes = 30
		}

		// 执行同步
		syncService := elasticsearch.NewSyncService(client)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		start := time.Now()
		if err := syncService.SyncTopicIncremental(ctx, minutes); err != nil {
			log.Fatalf("Incremental sync failed: %v", err)
		}

		elapsed := time.Since(start)
		fmt.Printf("\n✓ Incremental sync completed in %v\n", elapsed)
	},
}

// esStatusCmd 同步状态命令
var esStatusCmd = &cobra.Command{
	Use:   "sync-status",
	Short: "Check sync status between MySQL and Elasticsearch",
	Long:  "检查MySQL和Elasticsearch之间的同步状态",
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化数据库连接
		bootstrap.LoadDB()
		bootstrap.LoadLogger()

		// 创建ES客户端
		addresses := []string{"http://localhost:9200"}
		client, err := elasticsearch.NewClient(addresses)
		if err != nil {
			log.Fatalf("Failed to create ES client: %v", err)
		}

		// 获取同步状态
		syncService := elasticsearch.NewSyncService(client)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		status, err := syncService.GetSyncStatus(ctx)
		if err != nil {
			log.Fatalf("Failed to get sync status: %v", err)
		}

		fmt.Println("\n=== Elasticsearch Sync Status ===")
		fmt.Printf("MySQL Total Topics:    %v\n", status["mysql_total"])
		fmt.Printf("ES Indexed Topics:     %v\n", status["es_indexed"])
		fmt.Printf("Sync Status:           %v\n", status["sync_status"])
		if status["synced"].(bool) {
			fmt.Println("\n✓ All data is in sync!")
		} else {
			fmt.Println("\n⚠ Data is out of sync, consider running 'go run main.go es sync'")
		}
	},
}

// esReindexCmd 重建索引命令
var esReindexCmd = &cobra.Command{
	Use:   "reindex",
	Short: "Reindex all topics in Elasticsearch",
	Long:  "删除旧索引并重新创建，然后同步所有话题",
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化数据库连接
		bootstrap.LoadDB()
		bootstrap.LoadLogger()

		// 创建ES客户端
		addresses := []string{"http://localhost:9200"}
		client, err := elasticsearch.NewClient(addresses)
		if err != nil {
			log.Fatalf("Failed to create ES client: %v", err)
		}

		// 执行重建
		syncService := elasticsearch.NewSyncService(client)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		start := time.Now()
		if err := syncService.ReindexTopics(ctx, 1000); err != nil {
			log.Fatalf("Reindex failed: %v", err)
		}

		elapsed := time.Since(start)
		fmt.Printf("\n✓ Reindex completed in %v\n", elapsed)
	},
}

func init() {
	// 添加ES命令到根命令
	rootCmd.AddCommand(esCmd)

	// 添加子命令
	esCmd.AddCommand(esSyncCmd)
	esCmd.AddCommand(esIncrementalCmd)
	esCmd.AddCommand(esStatusCmd)
	esCmd.AddCommand(esReindexCmd)

	// 添加参数
	esIncrementalCmd.Flags().IntP("minutes", "m", 30, "Sync topics modified in the last N minutes")
}
