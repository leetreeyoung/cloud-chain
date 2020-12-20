package main

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/sinksmell/bee-crontab/models/worker"
)

func main() {
	var (
		err    error
		ctx    context.Context
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	// 初始化配置
	if err = worker.InitConfig(ctx, "./worker.yaml"); err != nil {
		goto ERR
	}

	// 启动日志协程
	if err = worker.InitLogger(ctx); err != nil {
		goto ERR
	}
	// 启动执行器
	if err = worker.InitExecutor(); err != nil {
		goto ERR
	}
	// 启动调度协程
	if err = worker.InitScheduler(); err != nil {
		log.Error("start scheduler err")
		goto ERR
	}
	// 启动任务管理器
	if err = worker.InitJobMgr(); err != nil {
		log.Error("start job manager err")
		goto ERR
	}

	// 启动服务注册
	if err = worker.InitRegister(); err != nil {
		log.Error("start init register err")
		goto ERR
	}

	// 启动监控服务
	if err = worker.InitPromMetrics(); err != nil {
		log.Error("start init prom metrics err")
		goto ERR
	}

	for {
		time.Sleep(100 * time.Millisecond)
	}

ERR:
	log.Errorf("worker start failed %v", err)

}
