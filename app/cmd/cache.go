package cmd

import (
	"GoHub-Service/pkg/cache"
	"GoHub-Service/pkg/console"

	"github.com/spf13/cobra"
)

var CmdCache = &cobra.Command{
	Use:   "cache",
	Short: "Cache management",
}

var CmdCacheClear = &cobra.Command{
	Use:   "clear",
	Short: "Clear cache",
	Run:   runCacheClear,
}

func init() {
	// 注册 cache 命令的子命令
	CmdCache.AddCommand(CmdCacheClear)
}

func runCacheClear(cmd *cobra.Command, args []string) {
	// 清空缓存
	cache.Flush()
	console.Success("Cache cleared.")
}
