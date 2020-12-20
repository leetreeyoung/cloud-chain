package worker

import (
	"context"

	"github.com/sinksmell/bee-crontab/models/common"
	log "github.com/sirupsen/logrus"
)

var (
	// BeeCronLogger 全局日志存储器单例
	BeeCronLogger *common.Logger
)

// InitLogger 初始化Logger的单例
func InitLogger(ctx context.Context) (err error) {

	BeeCronLogger, err = common.NewLogger(ctx, Conf.MongoURL)
	if err != nil {
		log.Errorf("init worker job logger err %w", err)
		return
	}

	// 启动日志存储协程
	go BeeCronLogger.WriteLoop()

	return
}
