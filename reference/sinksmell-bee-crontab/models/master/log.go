package master

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"sync"
	"time"

	"github.com/astaxie/beego"
	"github.com/sinksmell/bee-crontab/models/common"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	pool = sync.Pool{New: func() interface{} {
		mongoURL := beego.AppConfig.String("mongoURL")
		logger, err := common.NewLogger(context.Background(), mongoURL)
		if err != nil {
			panic(err)
		}
		return logger
	}}
)

// 创建一个新的 logger
func NewLogger() (logger *common.Logger) {
	return pool.Get().(*common.Logger)
}

// 回收 logger
func RecycleLogger(logger *common.Logger) {
	pool.Put(logger)
}

// ReadLog 读取任务的执行日志
func ReadLog(ctx context.Context, jobID string) (logs []*common.HTTPJobLog, err error) {

	var (
		execLog *common.JobExecLog
		httpLog *common.HTTPJobLog
		cursor  *mongo.Cursor
		findOps *options.FindOptions
		filter  *common.JobFilter
	)

	logger := NewLogger()
	defer RecycleLogger(logger)

	// 初始化返回结果 防止出现空指针
	logs = make([]*common.HTTPJobLog, 0)
	// 查找时的选项
	findOps = options.Find()
	findOps.SetLimit(20)
	// 设置过滤器即查找条件
	filter = &common.JobFilter{
		jobID,
	}
	findOps.SetSort(bson.M{
		"_id":-1,
	})
	if cursor, err = logger.LogCollection.Find(ctx, filter, findOps); err != nil {
		log.Errorf("find data err", err)
		return
	}
	// 延迟释放游标
	defer cursor.Close(ctx)
	// 遍历游标
	for cursor.Next(context.TODO()) {
		execLog = &common.JobExecLog{}
		err = cursor.Decode(execLog)
		if err != nil {
			log.Errorf("decode data err", err)
			continue
		}
		httpLog = toHTTPLog(execLog)
		logs = append(logs, httpLog)
	}

	return
}

//
func toHTTPLog(jobLog *common.JobExecLog) (log *common.HTTPJobLog) {

	// 构造http响应的log
	log = &common.HTTPJobLog{}
	log.JobName = jobLog.JobName
	log.Command = jobLog.Command
	log.Err = jobLog.Err
	log.Output = jobLog.Output
	// 时间戳转换为时间类型
	// time.Unix (seconds,nanoseconds)
	// 要么传入秒 要么传入纳秒
	// 由于之前获取的时毫秒级别的时间戳 这里将其转换为对应的毫秒
	log.StartTime = time.Unix(0, jobLog.StartTime*int64(time.Millisecond)).String()
	log.EndTime = time.Unix(0, jobLog.EndTime*int64(time.Millisecond)).String()
	log.PlanTime = time.Unix(0, jobLog.PlanTime*int64(time.Millisecond)).String()
	log.ScheduleTime = time.Unix(0, jobLog.ScheduleTime*int64(time.Millisecond)).String()

	return
}
