package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"webapp.demo/config"
	"webapp.demo/dao/mysql"
	"webapp.demo/dao/redis"
	"webapp.demo/logger"
	"webapp.demo/route"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("need config file. eg: app config.yaml")
		return
	}
	configFilePath := os.Args[1]
	// 1. 加载配置
	if err := config.Init(configFilePath); err != nil {
		fmt.Printf("init config failed. err: %v\n", err)
		return
	}
	// 2. 初始化日志
	if err := logger.Init(config.Conf.LogConfig); err != nil {
		fmt.Printf("init logger failed. err: %v\n", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("logger init success")

	// 3. 初始化Mysql连接
	if err := mysql.Init(config.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed. err: %v\n", err)
		return
	}
	defer mysql.Close()

	// 4. 初始化Redis连接
	if err := redis.Init(config.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed. err: %v\n", err)
		return
	}
	defer redis.Close()

	// 5. 注册路由
	r := route.Setup()
	// 6. 启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Conf.Port),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			zap.L().Fatal("listen", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅的关闭服务器，为关闭服务器设置一个5秒的超时
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server shutdown", zap.Error(err))
	}
	zap.L().Info("Server exit")
}
