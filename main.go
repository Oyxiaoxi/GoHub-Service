package main

import (
	"GoHub-Service/app/cmd"
	"GoHub-Service/app/cmd/make"
	"GoHub-Service/bootstrap"
	_ "GoHub-Service/config"
	"GoHub-Service/pkg/config"
	"GoHub-Service/pkg/console"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func main() {

	// 应用的主入口，默认调用 cmd.CmdServe 命令
	var rootCmd = &cobra.Command{
		Use:   "GoHub-Service",
		Short: "A simple forum project",
		Long:  `Default will run "serve" command, you can use "-h" flag to see all subcommands`,

		// PersistentPreRun 的所有子命令都会执行以下代码
		PersistentPreRun: func(command *cobra.Command, args []string) {

			// 配置初始化，依赖命令行 --env 参数
			config.InitConfig(cmd.Env)

			// 初始化日志（必须在其他组件之前）
			bootstrap.SetupLogger()

			// 初始化资源追踪器
			bootstrap.SetupTracker()

			// 初始化数据库
			bootstrap.SetupDB()

			// 初始化 Redis
			bootstrap.SetupRedis()

			// 初始化缓存
			bootstrap.SetupCache()

			// 启动资源泄漏定期报告（仅在 serve 命令时启动）
			if command.Name() == "serve" {
				// 5分钟未释放的资源视为可能泄漏，每1分钟检查一次
				bootstrap.StartTrackerReporting(5*time.Minute, 1*time.Minute)
			}
		},
	}

	// 注册子命令
	rootCmd.AddCommand(
		cmd.CmdServe,
		cmd.CmdKey,
		cmd.CmdPlay,
		make.CmdMake,
		cmd.CmdMigrate,
		cmd.CmdDBSeed,
		cmd.CmdCache,
		cmd.CmdSlowLog,
	)

	// 配置默认运行 Web 服务
	cmd.RegisterDefaultCmd(rootCmd, cmd.CmdServe)

	// 注册全局参数，--env
	cmd.RegisterGlobalFlags(rootCmd)

	// 执行主命令
	if err := rootCmd.Execute(); err != nil {
		console.Exit(fmt.Sprintf("Failed to run app with %v: %s", os.Args, err.Error()))
	}
}
