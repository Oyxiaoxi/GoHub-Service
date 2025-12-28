package cmd

import (
    "GoHub-Service/bootstrap"
    _ "GoHub-Service/docs" // Swagger 文档
    "GoHub-Service/pkg/config"
    "GoHub-Service/pkg/console"
    "GoHub-Service/pkg/logger"

    "github.com/gin-gonic/gin"
    ginSwagger "github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"
    "github.com/spf13/cobra"
)

// CmdServe represents the available web sub-command.
var CmdServe = &cobra.Command{
    Use:   "serve",
    Short: "Start web server",
    Run:   runWeb,
    Args:  cobra.NoArgs,
}

func runWeb(cmd *cobra.Command, args []string) {

    // 设置 gin 的运行模式，支持 debug, release, test
    // release 会屏蔽调试信息，官方建议生产环境中使用
    // 非 release 模式 gin 终端打印太多信息，干扰到我们程序中的 Log
    // 故此设置为 release，有特殊情况手动改为 debug 即可
    gin.SetMode(gin.ReleaseMode)

    // gin 实例
    router := gin.New()

    // 初始化路由绑定
    bootstrap.SetupRoute(router)

    // Swagger 文档路由
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // 运行服务器
    err := router.Run(":" + config.Get("app.port"))
    if err != nil {
        logger.ErrorString("CMD", "serve", err.Error())
        console.Exit("Unable to start server, error:" + err.Error())
    }
}
