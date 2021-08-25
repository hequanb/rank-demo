package main

import (
	"boframe/settings/mongoI"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"boframe/logger"
	"boframe/pkg/snowflake"
	"boframe/routers"
	"boframe/settings"
	"boframe/settings/redis"
	"go.uber.org/zap"
)

func main() {

	// 初始化配置文件
	if err := settings.Init(); err != nil {
		fmt.Printf("init setting failed: %v \n", err)
		return
	}

	// 初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init logger failed: %v \n", err)
		return
	}

	// 最后同步一下
	defer zap.L().Sync()
	zap.L().Debug("zap logger init success...")

	// 初始化Redis
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed: %v \n", err)
		return
	}
	defer redis.Close()

	// snowflake
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineId); err != nil {
		fmt.Printf("init snowflake failed: %v \n", err)
		return
	}

	// Mongo init
	if err := mongoI.Init(settings.Conf.MongoConfig); err != nil {
		fmt.Printf("init mongo failed: %v \n", err)
		return
	}
	defer mongoI.Close()

	// 注册路由
	r := routers.Setup()

	// 优雅停止
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: r,
	}


	// 异步启动，不然会阻塞
	go func() {
		zap.L().Debug("server starting...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen server failed", zap.Error(err))
		}
	}()

	// 等待中断信号到来，随后给5秒的反应时间，5秒后就不管了，强行关闭
	quit := make(chan os.Signal, 1)

	// 从下面的signal中，选择来转发给quit监听
	// syscall.SIGINT， kill -2, CTRL+C
	// syscall.SIGTERM, kill, 默认的关闭效果
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 因为用了带1个缓冲区的通道，这里不阻塞
	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// golang自带的优雅关闭服务，将未处理完的请求处理完再关闭，超过5s，强制关闭
	if err := server.Shutdown(ctx); err != nil {
		zap.L().Fatal("server shutdown err", zap.Error(err))
	}

	zap.L().Info("server exiting")
}
