// @title GoHub API
// @version 1.0
// @description GoHub 论坛服务 API 文档
// @termsOfService https://gohub.com/terms/

// @contact.name API Support
// @contact.url https://gohub.com/support
// @contact.email support@gohub.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
    "fmt"
    "GoHub-Service/app/cmd"
    "GoHub-Service/app/cmd/make"
    "GoHub-Service/bootstrap"
    _ "GoHub-Service/config"
    "GoHub-Service/pkg/config"
    "GoHub-Service/pkg/console"
    "os"

    "github.com/spf13/cobra"
)

func main() {

    // 应用的主入口，默认调用 cmd.CmdServe 命令
    var rootCmd = &cobra.Command{
        Use:   "GoHub-Service",
        Short: "A simple forum project",
        Long:  `Default will run "serve" command, you can use "-h" flag to see all subcommands`,

        // rootCmd 的所有子命令都会执行以下代码
        PersistentPreRun: func(command *cobra.Command, args []string) {

            // 配置初始化，依赖命令行 --env 参数
            config.InitConfig(cmd.Env)

            // 初始化 Logger
            bootstrap.SetupLogger()

            // 初始化数据库
            bootstrap.SetupDB()

            // 初始化 Redis
            bootstrap.SetupRedis()

            // 初始化缓存
            bootstrap.SetupCache()
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
