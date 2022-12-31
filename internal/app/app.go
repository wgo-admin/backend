package app

import (
  "context"
  "errors"
  "fmt"
  "github.com/wgo-admin/backend/pkg/validate"
  "net/http"
  "os"
  "os/signal"
  "syscall"
  "time"

  "github.com/gin-gonic/gin"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
  "github.com/wgo-admin/backend/internal/pkg/known"
  "github.com/wgo-admin/backend/internal/pkg/log"
  "github.com/wgo-admin/backend/pkg/token"
  "github.com/wgo-admin/backend/pkg/version/verflag"
)

var cfgFile string

func NewAppCommand() *cobra.Command {
  cmd := &cobra.Command{
    // 指定命令的名字，该名字会出现在帮助信息中
    Use: "wgo",
    // 命令的简短描述
    Short: "一个优质的Go后端项目",
    // 命令的详细描述
    Long: `基于 Gin 的后端服务项目`,

    // 命令出错时，不打印帮助信息。不需要打印帮助信息，设置为 true 可以保持命令出错时一眼就能看到错误信息
    SilenceUsage: true,
    // 指定调用 cmd.Execute() 时，执行的 Run 函数，函数执行失败会返回错误信息
    RunE: func(cmd *cobra.Command, args []string) error {
      // 如果 `--version=true`，则打印版本并退出
      verflag.PrintAndExitIfRequested()

      // 初始化日志
      log.Init(logConfig())
      // Sync 将缓存中的日志刷新到磁盘文件中
      defer log.Sync()

      return run()
    },
    // 这里设置命令运行时，不需要指定命令行参数
    Args: func(cmd *cobra.Command, args []string) error {
      for _, arg := range args {
        if len(arg) > 0 {
          return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
        }
      }

      return nil
    },
  }

  // 在执行命令前初始化配置
  cobra.OnInitialize(initConfig)

  // Cobra 支持持久性标志(PersistentFlag)，该标志可用于它所分配的命令以及该命令下的每个子命令
  cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "应用程序项目配置文件")

  // 添加 --version 标志
  verflag.AddFlags(cmd.PersistentFlags())

  return cmd
}

func run() error {
  // 初始化 store 层
  if err := initStore(); err != nil {
    return err
  }

  // 设置 token 包的签发密钥，用于 token 包的 token 签发和解析
  token.Init(viper.GetString("jwt-secret"), known.XUsernameKey, known.XRoleKey)

  // 初始化校验参数包
  validate.Init()

  // 设置 gin 模式
  gin.SetMode(viper.GetString("runmode"))

  // 创建 gin 引擎
  g := gin.New()

  // 注册公共中间件
  mws := []gin.HandlerFunc{gin.Recovery()}
  g.Use(mws...)

  // 注册路由
  if err := registryRoutes(g); err != nil {
    return err
  }

  // 创建并运行 http 服务器
  httpsrv := startInsecureServer(g)

  // TODO 创建并运行 https 服务器

  // 优雅关闭服务器
  quit := make(chan os.Signal, 1)
  signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
  <-quit
  log.Infow("Shutting down server ...")
  // 创建 ctx 用于通知服务器 goroutine, 它有 10 秒时间完成当前正在处理的请求
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  // 退出 http 服务器
  if err := httpsrv.Shutdown(ctx); err != nil {
    log.Errorw("Insecure Server forced to shutdown", "err", err)
    return err
  }

  log.Infow("Server exiting")

  return nil
}

// startInsecureServer 创建并运行 HTTP 服务器.
func startInsecureServer(g *gin.Engine) *http.Server {
  // 创建 HTTP Server 实例
  httpsrv := &http.Server{Addr: viper.GetString("addr"), Handler: g}

  // 运行 HTTP 服务器。在 goroutine 中启动服务器，它不会阻止下面的正常关闭处理流程
  // 打印一条日志，用来提示 HTTP 服务已经起来，方便排障
  log.Infow("Start to listening the incoming requests on http address", "addr", viper.GetString("addr"))
  go func() {
    if err := httpsrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
      log.Fatalw(err.Error())
    }
  }()

  return httpsrv
}
