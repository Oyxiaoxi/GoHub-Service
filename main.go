// Package main GoHub-Service API
// @title GoHub-Service API
// @version 1.0
// @description GoHub 社区论坛 API 文档
// @termsOfService https://github.com/Oyxiaoxi/GoHub-Service
// @contact.name API Support
// @contact.url https://github.com/Oyxiaoxi/GoHub-Service/issues
// @contact.email support@gohub.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:3000
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"GoHub-Service/app/cmd"
	"GoHub-Service/app/cmd/make"
	"GoHub-Service/bootstrap"
	_ "GoHub-Service/config"
	_ "GoHub-Service/docs" // Swagger 文档
	"GoHub-Service/pkg/appconfig"
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
				// 从配置读取阈值和检查间隔
				threshold := time.Duration(appconfig.GetResourceLeakThresholdMinutes()) * time.Minute
				interval := time.Duration(appconfig.GetCheckIntervalMinutes()) * time.Minute
				bootstrap.StartTrackerReporting(threshold, interval)
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
